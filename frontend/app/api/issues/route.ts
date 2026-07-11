import { NextResponse } from "next/server";
import { CAP_SITEVERIFY_URL } from "../../../lib/cap";

const GITHUB_TOKEN = process.env.GITHUB_TOKEN;
const CAP_SECRET = process.env.CAP_SECRET; // sunucu ortamında; git'te DEĞİL
const REPO_OWNER = "omerdduran";
const REPO_NAME = "turk-adfilter";

const MAX_TITLE = 200;
const MAX_BODY = 5000;

// Basit IP başına oran sınırı (turk-adfilter-prod tek instance → in-memory yeterli).
const RATE_LIMIT = 3; // saatte
const RATE_WINDOW_MS = 60 * 60 * 1000;
const hits = new Map<string, number[]>();

function isRateLimited(ip: string): boolean {
  const now = Date.now();
  const recent = (hits.get(ip) ?? []).filter((t) => now - t < RATE_WINDOW_MS);
  if (recent.length >= RATE_LIMIT) {
    hits.set(ip, recent);
    return true;
  }
  recent.push(now);
  hits.set(ip, recent);
  // Ara sıra eski girdileri temizle (bellek şişmesini önle).
  if (hits.size > 5000) {
    for (const [k, v] of hits) {
      const alive = v.filter((t) => now - t < RATE_WINDOW_MS);
      if (alive.length === 0) hits.delete(k);
      else hits.set(k, alive);
    }
  }
  return false;
}

// Cap token'ını sunucu tarafında doğrula.
// Dönüş: "ok" | "invalid" | "unreachable"
async function verifyCaptcha(token: string): Promise<"ok" | "invalid" | "unreachable"> {
  const controller = new AbortController();
  const timeout = setTimeout(() => controller.abort(), 5000);
  try {
    const res = await fetch(CAP_SITEVERIFY_URL, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ secret: CAP_SECRET, response: token }),
      signal: controller.signal,
    });
    if (!res.ok) return "unreachable";
    const data = (await res.json().catch(() => ({}))) as { success?: boolean };
    return data?.success === true ? "ok" : "invalid";
  } catch {
    return "unreachable"; // ağ hatası / timeout
  } finally {
    clearTimeout(timeout);
  }
}

export async function POST(request: Request) {
  if (!GITHUB_TOKEN) {
    return NextResponse.json({ message: "Sunucu yapılandırması eksik" }, { status: 500 });
  }

  const ip =
    request.headers.get("x-forwarded-for")?.split(",")[0]?.trim() ||
    request.headers.get("x-real-ip") ||
    "unknown";

  if (isRateLimited(ip)) {
    return NextResponse.json(
      { message: "Çok fazla istek gönderildi. Lütfen bir süre sonra tekrar deneyin." },
      { status: 429 }
    );
  }

  let payload: unknown;
  try {
    payload = await request.json();
  } catch {
    return NextResponse.json({ message: "Geçersiz istek" }, { status: 400 });
  }

  const { title, body, captchaToken } = (payload ?? {}) as {
    title?: unknown;
    body?: unknown;
    captchaToken?: unknown;
  };

  if (typeof title !== "string" || typeof body !== "string" || !title.trim() || !body.trim()) {
    return NextResponse.json({ message: "Başlık ve açıklama zorunludur" }, { status: 400 });
  }
  if (title.length > MAX_TITLE || body.length > MAX_BODY) {
    return NextResponse.json({ message: "Başlık veya açıklama çok uzun" }, { status: 400 });
  }

  // --- CAPTCHA doğrulama (fail-safe) ---
  if (!CAP_SECRET) {
    return NextResponse.json({ message: "Doğrulama yapılandırılmamış" }, { status: 500 });
  }
  if (typeof captchaToken !== "string" || !captchaToken) {
    return NextResponse.json({ message: "Güvenlik doğrulaması gerekli" }, { status: 400 });
  }
  const verdict = await verifyCaptcha(captchaToken);
  if (verdict === "invalid") {
    return NextResponse.json(
      { message: "Güvenlik doğrulaması geçersiz. Lütfen tekrar deneyin." },
      { status: 400 }
    );
  }
  // verdict === "unreachable": Cap servisi anlık erişilemez — form kilitlenmesin
  // diye oran sınırına güvenerek geçir, ama logla.
  if (verdict === "unreachable") {
    console.warn("[issues] Cap siteverify erişilemedi; rate-limit ile fail-open (ip=%s)", ip);
  }

  try {
    const response = await fetch(
      `https://api.github.com/repos/${REPO_OWNER}/${REPO_NAME}/issues`,
      {
        method: "POST",
        headers: {
          Authorization: `Bearer ${GITHUB_TOKEN}`,
          Accept: "application/vnd.github+json",
          "Content-Type": "application/json",
        },
        // labels istemciden ALINMAZ — sabit atanır.
        body: JSON.stringify({ title, body, labels: ["user-feedback"] }),
      }
    );

    if (!response.ok) {
      const err = await response.json().catch(() => ({}));
      throw new Error((err as { message?: string })?.message || "GitHub hatası");
    }

    const data = (await response.json()) as { html_url?: string; number?: number };
    return NextResponse.json({ url: data.html_url, number: data.number });
  } catch (error) {
    console.error("[issues] oluşturma hatası:", error);
    return NextResponse.json(
      { message: "Gönderilemedi. Lütfen daha sonra tekrar deneyin." },
      { status: 502 }
    );
  }
}
