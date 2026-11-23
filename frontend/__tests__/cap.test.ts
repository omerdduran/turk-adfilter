import { describe, expect, it } from "vitest";
import { CAP_INSTANCE, CAP_SITE_KEY, CAP_SITEVERIFY_URL, CAP_WIDGET_ENDPOINT } from "../lib/cap";

describe("cap configuration", () => {
  it("builds widget and siteverify endpoints from the public instance", () => {
    expect(CAP_WIDGET_ENDPOINT).toBe(`${CAP_INSTANCE}/${CAP_SITE_KEY}/`);
    expect(CAP_SITEVERIFY_URL).toBe(`${CAP_INSTANCE}/${CAP_SITE_KEY}/siteverify`);
  });

  it("uses https for the captcha instance", () => {
    expect(CAP_INSTANCE.startsWith("https://")).toBe(true);
  });

  it("never exposes a secret key", () => {
    expect(JSON.stringify({ CAP_INSTANCE, CAP_SITE_KEY, CAP_SITEVERIFY_URL })).not.toMatch(
      /secret/i,
    );
  });
});
