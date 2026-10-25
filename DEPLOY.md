# Deploy

Üretim ortamı **elle** dağıtılır; hiçbir CI sunucusu bu depoya bağlı değil
(`.woodpecker.yml` bu yüzden kaldırıldı — kimseye çalışıyormuş gibi görünmesin).

## Adımlar

Sunucuda `/opt/turk-adfilter/` altında iki şey var: `prod/` (docker-compose + `.env`) ve
`_src/frontend/` (derlenecek ağaç).

```bash
# 1. Yerelde main'e geç ve ağacı sunucuya eşitle
cd frontend
rsync -a --delete --exclude node_modules --exclude .next --exclude .source \
  ./ main:/opt/turk-adfilter/_src/frontend/

# 2. Sunucuda derle + yeniden başlat
ssh main 'CI_COMMIT_BRANCH=main CI_WORKSPACE=/opt/turk-adfilter/_src \
  sh /opt/turk-adfilter/deploy.sh'
```

`deploy.sh` `frontend/`'i `turk-adfilter-web:release` imajı olarak derler, `prod/`
dizininde `docker compose up -d` çalıştırır ve `/opt/turk-adfilter/deploy.log`'a
START/DONE yazar. Yalnızca `main` kabul edilir.

## Doğrulama

```bash
curl -sSL -o /dev/null -w '%{http_code}\n' https://reklamsiz-turkiye.com/
ssh main 'docker ps --format "{{.Names}}\t{{.Status}}" | grep turk-adfilter-prod'
```

## Geri alma

İmaj `:release` etiketiyle üzerine yazıldığı için eski sürüm diskte durmaz. Geri
dönmek için ilgili sürüm etiketinden yeniden derleyip dağıtın:

```bash
git checkout v1.9.7        # ya da sorunlu olmayan son sürüm
rsync -a ...               # yukarıdaki 1. adım
ssh main 'CI_COMMIT_BRANCH=main CI_WORKSPACE=/opt/turk-adfilter/_src sh /opt/turk-adfilter/deploy.sh'
```

## Kaldırılan ortamlar

- **preview.reklamsiz-turkiye.com** (2026-10-25): staging ortamı, Caddy site bloğu,
  `preview/` compose dizini, `:preview` imajı ve DNS kaydı kaldırıldı. Prod elle
  dağıtıldığı için ayrı bir staging ortamı tutulmuyor.

## Sırlar

Yalnızca sunucuda, `/opt/turk-adfilter/prod/.env` içinde: `GITHUB_TOKEN`, `CAP_SECRET`,
`NEXTAUTH_URL`. Depoda sır yok; `.gitignore` yerel `.env` dosyalarını da dışlar.
