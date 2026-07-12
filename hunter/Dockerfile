# hunter — turk-adfilter otomatik domain bulucu bot
# Çok aşamalı: statik binary → distroless (CA + tzdata hazır, ~15MB, nonroot).
FROM golang:1.25-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /hunter ./cmd/hunter
# /data'yı nonroot (65532) sahipliğiyle oluştur — boş named/anonymous volume bu izni
# miras alır, böylece distroless nonroot kullanıcı SQLite'ı yazabilir.
RUN mkdir -p /data && chown 65532:65532 /data

FROM gcr.io/distroless/static-debian12:nonroot
LABEL org.opencontainers.image.source="https://github.com/omerdduran/turk-adfilter"
LABEL org.opencontainers.image.description="turk-adfilter otomatik domain bulucu bot"
COPY --from=build /hunter /hunter
COPY --from=build --chown=65532:65532 /data /data
VOLUME ["/data"]
# Dead-man's-switch: son başarılı çalışma çok eskiyse unhealthy (bot asılmış/ölmüş).
HEALTHCHECK --interval=1h --timeout=10s --start-period=5m --retries=2 \
	CMD ["/hunter", "healthcheck"]
ENTRYPOINT ["/hunter"]
CMD ["serve"]
