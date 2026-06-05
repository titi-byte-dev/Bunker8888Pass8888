# AegisPass — Site institucional

Site publico estatico (`aegispass.com`), **independente** da app
(`app.aegispass.com`). Apresenta o ecossistema; nao tem login nem sessao.

Spec: [`docs/roadmap/09-design/site-institucional-spec.md`](../docs/roadmap/09-design/site-institucional-spec.md)

## Principios

- **Separado da app.** Codigo, build e deploy proprios. So partilha tokens.
- **Estatico.** SvelteKit `adapter-static`, prerender total, SEO-friendly.
- **Fonte unica de marca.** `tokens.css` da app e copiado por `sync-tokens.mjs`.
- **Fonte unica de copy.** Texto em `src/lib/content.ts` (prepara i18n).

## Comandos

```bash
npm install
npm run dev      # copia tokens + arranca dev
npm run build    # site estatico em build/
npm run check    # svelte-check
npm test         # vitest (guarda a narrativa)
```

## Fronteira

```
frontend/src/lib/design/tokens.css   (FONTE de marca)
        | sync-tokens.mjs (predev/prebuild)
        v
site/src/lib/tokens.generated.css    (artefacto, gitignored)

aegispass.com  --[Entrar]-->  app.aegispass.com
```
