import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";

interface MockResponse {
  ok: boolean;
  status?: number;
  json?: unknown;
}

function mockFetch(responses: MockResponse[]) {
  const fetchMock = vi.fn(async () => {
    const next = responses.shift();
    if (!next) throw new Error("beklenmeyen fetch çağrısı");
    return {
      ok: next.ok,
      status: next.status ?? (next.ok ? 200 : 500),
      json: async () => next.json ?? {},
    } as Response;
  });
  vi.stubGlobal("fetch", fetchMock);
  return fetchMock;
}

async function importRoute({ token = "gh-token", secret = "cap-secret" } = {}) {
  vi.resetModules();
  vi.stubEnv("GITHUB_TOKEN", token);
  vi.stubEnv("CAP_SECRET", secret);
  return import("../app/api/issues/route");
}

function request(body: BodyInit | null, ip = "1.2.3.4") {
  return new Request("http://localhost/api/issues", {
    method: "POST",
    headers: { "content-type": "application/json", "x-forwarded-for": ip },
    body,
  });
}

const validBody = {
  title: "Reklam bildirimi",
  description: "Site reklam gösteriyor",
  priority: "high",
  captchaToken: "token",
};

beforeEach(() => {
  vi.stubEnv("GITHUB_TOKEN", undefined);
  vi.stubEnv("CAP_SECRET", undefined);
});

afterEach(() => {
  vi.unstubAllEnvs();
  vi.unstubAllGlobals();
  vi.restoreAllMocks();
});

describe("POST /api/issues", () => {
  it("returns 500 when GITHUB_TOKEN is missing", async () => {
    const { POST } = await importRoute({ token: "" });
    const response = await POST(request(JSON.stringify(validBody)));
    expect(response.status).toBe(500);
  });

  it("returns 500 when CAP_SECRET is missing", async () => {
    const { POST } = await importRoute({ secret: "" });
    const response = await POST(request(JSON.stringify(validBody)));
    expect(response.status).toBe(500);
  });

  it("rejects invalid JSON", async () => {
    const { POST } = await importRoute();
    const response = await POST(request("{bozuk"));
    expect(response.status).toBe(400);
  });

  it("requires title and description", async () => {
    const { POST } = await importRoute();
    const response = await POST(request(JSON.stringify({ title: " ", description: "" })));
    expect(response.status).toBe(400);
  });

  it("rejects oversized fields", async () => {
    const { POST } = await importRoute();
    const response = await POST(request(JSON.stringify({ ...validBody, title: "x".repeat(201) })));
    expect(response.status).toBe(400);
  });

  it("rejects an invalid captcha token", async () => {
    mockFetch([{ ok: true, json: { success: false } }]);
    const { POST } = await importRoute();
    const response = await POST(request(JSON.stringify(validBody)));
    expect(response.status).toBe(400);
  });

  it("creates the issue with server-side labels and assignee", async () => {
    const fetchMock = mockFetch([
      { ok: true, json: { success: true } },
      { ok: true, status: 201, json: { html_url: "https://github.com/x/y/issues/1", number: 1 } },
    ]);
    const { POST } = await importRoute();

    const response = await POST(request(JSON.stringify(validBody)));
    expect(response.status).toBe(200);
    await expect(response.json()).resolves.toEqual({
      url: "https://github.com/x/y/issues/1",
      number: 1,
    });

    const githubCall = fetchMock.mock.calls[1] as unknown as [string, RequestInit];
    expect(githubCall[0]).toContain("/repos/omerdduran/turk-adfilter/issues");
    const payload = JSON.parse(githubCall[1].body as string);
    expect(payload.labels).toEqual(["user-feedback", "priority: high"]);
    expect(payload.assignees).toEqual(["omerdduran"]);
  });

  it("rate limits after three requests from the same IP", async () => {
    mockFetch([
      { ok: true, json: { success: false } },
      { ok: true, json: { success: false } },
      { ok: true, json: { success: false } },
    ]);
    const { POST } = await importRoute();

    for (let i = 0; i < 3; i += 1) {
      const response = await POST(request(JSON.stringify(validBody), "9.9.9.9"));
      expect(response.status).toBe(400);
    }

    const response = await POST(request(JSON.stringify(validBody), "9.9.9.9"));
    expect(response.status).toBe(429);
  });
});
