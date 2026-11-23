import { createMDX } from "fumadocs-mdx/next";

const withMDX = createMDX();

/** @type {import('next').NextConfig} */
const config = {
  reactStrictMode: true,
  // Next 16 infers the workspace root from lockfiles. Pin it to the app so
  // standalone output stays at .next/standalone/server.js (the Dockerfile path).
  turbopack: {
    root: process.cwd(),
  },
  // Self-host (Docker) için minimal Node.js sunucusu üretir: .next/standalone
  // içinde server.js + gerekli node_modules bundle'lanır. Vercel bunu yok sayar.
  output: "standalone",
  // Güvenlik başlıkları (clickjacking / MIME-sniffing / referrer sızıntısı / CSP).
  //
  // CSP notları:
  // - Next inline hydration script'leri nonce olmadan çalıştığı için script-src'te
  //   'unsafe-inline' gerekli. Nonce, sayfaları dinamik render'a zorlayacağı için
  //   static/SSG çıktısını korumak adına tercih edilmedi.
  // - Cap widget WebAssembly ve blob worker kullanıyor: 'wasm-unsafe-eval' ve
  //   worker-src blob: yalnızca bu ihtiyaç için açıldı (JS eval'e izin verilmiyor).
  // - Cap widget'ı cdn.jsdelivr.net'ten, doğrulama isteği captcha alt alan adına gider.
  async headers() {
    return [
      {
        source: "/:path*",
        headers: [
          {
            key: "Content-Security-Policy",
            value: [
              "default-src 'self'",
              "script-src 'self' 'unsafe-inline' 'wasm-unsafe-eval' https://cdn.jsdelivr.net",
              "style-src 'self' 'unsafe-inline'",
              "img-src 'self' data:",
              "font-src 'self' data:",
              "connect-src 'self' https://captcha.reklamsiz-turkiye.com",
              "worker-src 'self' blob:",
              "frame-src 'self' blob: data:",
              "object-src 'none'",
              "base-uri 'self'",
              "form-action 'self'",
              "frame-ancestors 'self'",
              "upgrade-insecure-requests",
            ].join("; "),
          },
          { key: "X-Frame-Options", value: "SAMEORIGIN" },
          { key: "X-Content-Type-Options", value: "nosniff" },
          { key: "Referrer-Policy", value: "strict-origin-when-cross-origin" },
          {
            key: "Permissions-Policy",
            value: "camera=(), microphone=(), geolocation=(), interest-cohort=()",
          },
        ],
      },
    ];
  },
};

export default withMDX(config);
