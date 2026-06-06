# Deploy do site institucional (WEB-004)

> **aegispass.com** — Astro estático em `website/dist/`. Independente da app
> (`app.aegispass.com`).

## Build

```bash
cd website
npm ci
PUBLIC_APP_URL=https://app.aegispass.com npm run build
```

Output: `website/dist/` (4 páginas: `/`, `/fr/`, `/es/`, `/de/`).

## nginx (VPS)

1. Copiar `website/deploy/nginx-aegispass.conf.example` para
   `/etc/nginx/sites-available/aegispass.com`.
2. Ajustar `root` para o path do `dist/` no servidor.
3. Certificado TLS (Let's Encrypt / certbot).
4. `nginx -t && systemctl reload nginx`.

## Domínios

| Host | Destino |
|------|---------|
| `aegispass.com` | `website/dist/` |
| `www.aegispass.com` | redirect → apex |
| `app.aegispass.com` | frontend SvelteKit (Docker) |

## CI

O job `Website (Astro)` no GitHub Actions valida `npm run build` + testes i18n
em cada push/PR.

## Variáveis

| Variável | Uso |
|----------|-----|
| `PUBLIC_APP_URL` | URL do CTA «Entrar» (build-time) |
