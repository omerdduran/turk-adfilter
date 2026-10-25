"""Bounded DNS probing shared by filter_stats.py and domain_check.py.

`socket.gethostbyname` has no timeout knob, and every blocked or dead domain waits
for the system resolver timeout. Probing the 6 400+ list domains serially measured
~8 minutes locally (72 ms per lookup) and stalled a GitHub runner past the useful
part of the job, so lookups run in daemon threads and whatever has not answered by
the deadline is reported as unknown instead of hanging the process.

The deadline is deliberately part of the API: a partial answer plus a count of
unresolved domains is honest, a workflow that never finishes is not.
"""

import os
import queue
import socket
import threading
import time

DEFAULT_WORKERS = 256
DEFAULT_TIMEOUT = 180.0


def workers():
    return int(os.environ.get("DNS_PROBE_WORKERS", DEFAULT_WORKERS))


def deadline_seconds():
    return float(os.environ.get("DNS_PROBE_TIMEOUT", DEFAULT_TIMEOUT))


def resolve_many(domains, workers=None, timeout=None):
    """Resolve domains in parallel.

    Returns {domain: ip} for every domain that answered before the deadline, with
    ip = None for a domain that resolved to nothing (NXDOMAIN, SERVFAIL). Domains
    missing from the dict did not answer in time.
    """
    domains = list(domains)
    if not domains:
        return {}

    workers = max(1, min(workers or workers_default(), len(domains)))
    timeout = timeout or deadline_seconds()
    end = time.monotonic() + timeout

    results = {}
    lock = threading.Lock()
    pending = queue.Queue()
    for domain in domains:
        pending.put(domain)

    def worker():
        while True:
            try:
                domain = pending.get_nowait()
            except queue.Empty:
                return
            try:
                ip = socket.gethostbyname(domain)
            except Exception:
                ip = None
            with lock:
                results[domain] = ip

    threads = [
        threading.Thread(target=worker, name=f"dns-probe-{i}", daemon=True)
        for i in range(workers)
    ]
    for thread in threads:
        thread.start()
    for thread in threads:
        remaining = end - time.monotonic()
        if remaining <= 0:
            break
        thread.join(remaining)

    with lock:
        return dict(results)


def workers_default():
    return workers()
