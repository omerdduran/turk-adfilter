# syntax=docker/dockerfile:1.7
# Turk-AdFilter web sitesi (Next.js 16 + fumadocs) — self-host için standalone imaj.
# Site dinamiktir (app/api/*), bu yüzden static export değil Node.js runtime kullanılır.
# Domain (metadataBase) kodda sabit olduğu için build-arg gerekmez; GITHUB_TOKEN /
# NEXTAUTH_URL çalışma anında docker compose .env'den gelir.

# ── deps: bağımlılıklar (postinstall'i atla — fumadocs-mdx içeriği build'de üretilir) ──
FROM node:22-alpine AS deps
WORKDIR /app
COPY package.json package-lock.json ./
RUN --mount=type=cache,target=/root/.npm \
    npm ci --no-audit --no-fund --ignore-scripts

# ── builder: .source üretimi + next build (standalone çıktı) ──
FROM node:22-alpine AS builder
WORKDIR /app
ENV NEXT_TELEMETRY_DISABLED=1
COPY --from=deps /app/node_modules ./node_modules
COPY . .
RUN npm run build

# ── runner: minimal standalone sunucu, non-root ──
FROM node:22-alpine AS runner
WORKDIR /app
RUN apk add --no-cache tini wget && \
    addgroup -S app && adduser -S app -G app
ENV NODE_ENV=production \
    NEXT_TELEMETRY_DISABLED=1 \
    HOSTNAME=0.0.0.0 \
    PORT=3000
COPY --from=builder /app/public ./public
COPY --from=builder /app/.next/standalone ./
COPY --from=builder /app/.next/static ./.next/static
USER app
EXPOSE 3000
HEALTHCHECK --interval=10s --timeout=3s --start-period=15s --retries=5 \
  CMD wget -qO- http://127.0.0.1:3000/ >/dev/null 2>&1 || exit 1
ENTRYPOINT ["/sbin/tini", "--"]
CMD ["node", "server.js"]
