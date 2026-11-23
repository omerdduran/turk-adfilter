import { describe, expect, it } from "vitest";
import { createRateLimiter } from "../lib/rate-limit";

describe("createRateLimiter", () => {
  it("allows requests up to the limit and rejects the next one", () => {
    const now = 1_000;
    const limiter = createRateLimiter({ limit: 2, windowMs: 1_000, now: () => now });

    expect(limiter.isRateLimited("ip")).toBe(false);
    expect(limiter.isRateLimited("ip")).toBe(false);
    expect(limiter.isRateLimited("ip")).toBe(true);
    expect(limiter.isRateLimited("ip")).toBe(true);
  });

  it("accepts requests again after the window expires", () => {
    let now = 0;
    const limiter = createRateLimiter({ limit: 1, windowMs: 1_000, now: () => now });

    expect(limiter.isRateLimited("ip")).toBe(false);
    expect(limiter.isRateLimited("ip")).toBe(true);

    now = 1_000;
    expect(limiter.isRateLimited("ip")).toBe(false);
  });

  it("keeps separate counters per key", () => {
    const limiter = createRateLimiter({ limit: 1, windowMs: 1_000, now: () => 0 });

    expect(limiter.isRateLimited("a")).toBe(false);
    expect(limiter.isRateLimited("b")).toBe(false);
    expect(limiter.isRateLimited("a")).toBe(true);
  });

  it("does not count timestamps outside the window", () => {
    let now = 0;
    const limiter = createRateLimiter({ limit: 2, windowMs: 100, now: () => now });

    expect(limiter.isRateLimited("ip")).toBe(false);
    now = 50;
    expect(limiter.isRateLimited("ip")).toBe(false);
    now = 120;
    expect(limiter.isRateLimited("ip")).toBe(false);
    expect(limiter.isRateLimited("ip")).toBe(true);
  });
});
