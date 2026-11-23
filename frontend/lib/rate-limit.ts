export interface RateLimiterOptions {
  limit: number;
  windowMs: number;
  maxKeys?: number;
  now?: () => number;
}

export interface RateLimiter {
  isRateLimited(key: string): boolean;
}

export function createRateLimiter({
  limit,
  windowMs,
  maxKeys = 5000,
  now = Date.now,
}: RateLimiterOptions): RateLimiter {
  const hits = new Map<string, number[]>();

  return {
    isRateLimited(key: string): boolean {
      const current = now();
      const recent = (hits.get(key) ?? []).filter((timestamp) => current - timestamp < windowMs);

      if (recent.length >= limit) {
        hits.set(key, recent);
        return true;
      }

      recent.push(current);
      hits.set(key, recent);

      if (hits.size > maxKeys) {
        for (const [storedKey, timestamps] of hits) {
          const alive = timestamps.filter((timestamp) => current - timestamp < windowMs);
          if (alive.length === 0) hits.delete(storedKey);
          else hits.set(storedKey, alive);
        }
      }

      return false;
    },
  };
}
