# Site institucional AegisPass

One-pager estático para **aegispass.com** — independente da app SvelteKit
(`app.aegispass.com`). Partilha apenas `frontend/src/lib/design/tokens.css`.

## Idiomas

| Locale | URL |
|--------|-----|
| Português (PT) | `/` |
| Français | `/fr/` |
| Español | `/es/` |
| Deutsch | `/de/` |

Copy em `src/i18n/locales/*.ts` — adicionar idioma = novo ficheiro + entrada em
`astro.config.mjs` e `src/config.ts`.

## Desenvolvimento

```bash
cd website
npm install
npm run dev
```

Abre `http://localhost:4321`. Para apontar o CTA «Entrar» ao frontend local:

```bash
PUBLIC_APP_URL=http://localhost:5173 npm run dev
```

## Build

```bash
npm run build
npm run preview
```

Output em `website/dist/` — servir via CDN ou nginx estático.

Deploy em produção: [`docs/roadmap/10-production/website-deploy.md`](../docs/roadmap/10-production/website-deploy.md)
e [`deploy/nginx-aegispass.conf.example`](deploy/nginx-aegispass.conf.example).

## PRD

Ver [`docs/roadmap/09-design/site-institucional-spec.md`](../docs/roadmap/09-design/site-institucional-spec.md).
