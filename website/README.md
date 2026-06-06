# Site institucional AegisPass

One-pager **multi-página** para **aegispass.com** — Astro estático, independente da app.
Partilha `frontend/src/lib/design/tokens.css`. Copy em **PT, FR, ES, DE**.

## Estrutura (estilo Linear)

| Rota (PT) | Página |
|-----------|--------|
| `/` | Campanha + grelha de produtos |
| `/platform/` | Zero-knowledge + ecossistema em camadas |
| `/products/vault/` | Cofre |
| `/products/security/` | Segurança |
| `/products/team/` | Equipa |
| `/products/workspace/` | Workspace BYOD |
| `/partners/` | White-label / parceiros |

Locales: `/fr/…`, `/es/…`, `/de/…` (selector mantém o path actual).

Estados por página: **Disponível** · **Pré-visualização** · **Em construção** (banner + badges).

## Desenvolvimento

```bash
cd website
npm install
npm run dev
```

`npm run scaffold:pages` regenera wrappers finos em `src/pages/` a partir dos templates.

CTA «Entrar»:

```bash
PUBLIC_APP_URL=http://localhost:5173 npm run dev
```

## Build

```bash
npm run build   # 28 páginas estáticas
npm run test    # paridade i18n
```

Deploy: [`docs/roadmap/10-production/website-deploy.md`](../docs/roadmap/10-production/website-deploy.md)

## Adicionar copy

1. Editar `src/i18n/locales/{pt,fr,es,de}.ts` (mesmas chaves — testes validam).
2. Templates em `src/templates/` · componentes em `src/components/`.

PRD: [`docs/roadmap/09-design/site-institucional-spec.md`](../docs/roadmap/09-design/site-institucional-spec.md)
