import { createMDX } from 'fumadocs-mdx/next';

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
  output: 'standalone',
  // Temel güvenlik başlıkları (clickjacking / MIME-sniffing / referrer sızıntısı).
  // Not: CSP bilinçli olarak eklenmedi — Cap widget (jsdelivr), inline JSON-LD ve
  // font kaynaklarıyla uyumu ayrı test gerektirir; sonra dar bir CSP eklenecek.
  async headers() {
    return [
      {
        source: '/:path*',
        headers: [
          { key: 'X-Frame-Options', value: 'SAMEORIGIN' },
          { key: 'X-Content-Type-Options', value: 'nosniff' },
          { key: 'Referrer-Policy', value: 'strict-origin-when-cross-origin' },
          {
            key: 'Permissions-Policy',
            value: 'camera=(), microphone=(), geolocation=(), interest-cohort=()',
          },
        ],
      },
    ];
  },
};

export default withMDX(config);
