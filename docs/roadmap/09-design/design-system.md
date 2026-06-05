# Design System — AegisPass

> Plano operacional para evoluir de **tokens + shell** para uma app **clean,
> intuitiva e futurista** (sem cyberpunk), com hierarquia de navegação clara,
> várias paletas de tema e componentes reutilizáveis.
>
> Complementa [`product-ui-vision.md`](product-ui-vision.md) (North Star) e as
> tasks `UI-*` em [`../05-tasks/backlog.md`](../05-tasks/backlog.md).

---

## 1. Estado actual (Jun 2026)

| Área | Ficheiros / rotas | Estado | Gap |
|---|---|---|---|
| Tokens | `frontend/src/lib/design/tokens.css` (UI-001) | 🟢 | Só `light` / `dark` / `system` |
| Tema | `frontend/src/lib/design/theme.ts` | 🟢 | Sem paletas múltiplas |
| Shell | `lib/shell/AppSidebar`, `AppTabBar`, `(app)/+layout.svelte` | 🟢 | Navegação **plana** (~15 links) |
| Motion | `lib/motion/*` (UI-005) | 🟢 | — |
| Command palette | `lib/shell/CommandPalette.svelte` (UI-006) | 🟢 | Comandos flat, sem árvore |
| Catálogo dev | `/dev/components` (UI-010) | 🟢 | Botões/forms **inline**, não importáveis |
| Componentes UI | — | ❌ | **Não existe** `lib/ui/` |
| Breadcrumbs | — | ❌ | Só `.eyebrow` com IDs de task |
| Páginas | ~42 rotas em `(app)/` | 🟡 | 3 padrões visuais inconsistentes |

> 💡 **Conceito — Design debt:** quando tokens existem mas cada página redefine
> `.page-head`, `.panel` e `.eyebrow` localmente, o sistema **parece** maduro mas
> **comporta-se** como protótipo. A dívida manifesta-se em inconsistência UX e
> custo de manutenção.

### 1.1 Padrões visuais hoje (inconsistentes)

1. **Hub** — lista de links descritivos (`/work`, `/security`): `h1` + `.lead` + `.links`
2. **Dados densos** — painéis e formulários (`/fin`, `/crm`, `/hr`): `.page-head` + `.eyebrow` + `.panel`
3. **Cofre** — padrão próprio (`/vault`): `.page-header`, lista com motion, empty state

Os IDs de task nos eyebrows (`FIN-001`, `GOOGLE-001`) são úteis para **dev**; em
produção o utilizador deve ver **breadcrumbs** com nomes de produto.

---

## 2. North Star — futurista sem cyberpunk

Alinhado a [`product-ui-vision.md` §2](product-ui-vision.md): **Proton × Linear × Apple Wallet**.

| Princípio | Tradução visual |
|---|---|
| Confiança visível | Estados de cifragem, sessão, turno legíveis — nunca alarmistas |
| Calma operacional | Animações ≤400ms, poucos modais, zero urgência falsa |
| Densidade inteligente | Admin denso; funcionário espaçado (menor privilégio na UI) |
| Futuro = fluidez | 60fps, profundidade subtíl, feedback imediato |

**Futurista** neste projeto significa:

- **Fluidez** — cross-fade entre páginas (`ShellPageMotion`), stagger em listas
- **Profundidade** — camadas `base → elevated → surface` + `backdrop-blur` leve
- **Luz contida** — accent `#4DA3FF`; gradientes só em hero/empty (opacidade ≤15%)
- **Micro-feedback** — copy com countdown, toasts, skeletons

**Evitar:** neon, glitch, sombras pesadas, animações agressivas, modais em cascata.

---

## 3. Arquitectura de informação (sidebar + breadcrumbs)

### 3.1 Problema actual

`lib/shell/nav.ts` expõe subpáginas ao **mesmo nível** que módulos top-level
(ex.: `Fiscal` e `Faturas` como irmãos de `Cofre`). Isto contradiz
[`product-ui-vision.md` §4.1](product-ui-vision.md).

### 3.2 Modelo alvo — árvore de navegação

> 💡 **Conceito — Hub:** página intermédia que agrupa sub-rotas relacionadas
> (`/fin`, `/work`, `/security`). A sidebar mostra o **módulo**; os **filhos**
> expandem-se quando activos.

```
COFRE          → /vault
SEGURANÇA      → /security (hub)
  ├ Saúde de segurança    → /security/hygiene
  ├ Dispositivos e sessões → /security/devices
  ├ Sentinel Mode         → /security/sentinel
  ├ Acesso de emergência  → /security/emergency
  └ Auditoria Guardião    → /security/guardian
TRABALHO       → /work (hub)
  ├ Turnos e geofence     → /work/shifts
  ├ Browser sandbox       → /work/sandbox
  ├ CLI mTLS              → /work/cli
  ├ Inventário            → /work/inventory
  └ Google Workspace      → /work/google
EQUIPA         → /team (hub)
  ├ Shared vaults         → /team/vaults
  ├ Notas partilhadas     → /team/notes
  └ Secret links          → /team/links
RH             → /hr (hub)
  ├ Fichas e contratos    → /hr
  ├ Onboarding            → /hr/onboarding
  ├ Recrutamento          → /hr/recruitment
  └ Conformidade RGPD     → /hr/compliance
MAIL           → /mail
FINANÇAS       → /fin (hub)          ← não «Custos» solto na sidebar
  ├ Custos SaaS           → /fin
  ├ Fiscal                → /fin/fiscal
  ├ Faturas               → /fin/invoices
  ├ Comissões             → /fin/commissions
  └ Reconciliação bancária → /fin/banking
CRM            → /crm
ADMIN          → /admin (hub)
  ├ Utilizadores          → /admin/users
  └ Audit log             → /admin/audit
DEFINIÇÕES     → /settings
```

### 3.3 Tab bar mobile (máx. 5 itens)

Manter **hubs** apenas: Cofre · Segurança · Trabalho · Definições (+ opcional CRM/Fin para admin).

Sub-páginas: hub do módulo, breadcrumbs ou ⌘K.

### 3.4 Fonte de verdade — `ROUTE_TREE`

Ficheiro proposto: `frontend/src/lib/shell/routes.ts`

Responsabilidades únicas:

| Consumidor | Uso |
|---|---|
| `AppSidebar.svelte` | Secções colapsáveis + filho activo expande pai |
| `Breadcrumbs.svelte` | Trilho `Finanças › Fiscal` derivado do pathname |
| `CommandPalette` | Comandos hierárquicos (`Finanças › Fiscal`) |
| `+page.ts` (override) | Títulos dinâmicos (`vault/[id]` → nome do item) |

```ts
// Esboço — não implementado ainda (UI-011)
type RouteNode = {
  label: string;
  href?: string;
  taskId?: string;       // só eyebrow em DEV
  children?: Record<string, RouteNode>;
};
```

### 3.5 Breadcrumbs

Componente `Breadcrumbs.svelte` dentro de `PageShell`:

```
Finanças › Fiscal
Trabalho › Google Workspace
Cofre › Netflix › Editar
```

Regras:

- **Produção:** nomes de produto em PT-PT
- **Dev:** eyebrow opcional com `taskId` (`FIN-005`) via `import.meta.env.DEV`
- Último segmento não é link; anteriores são `<a href>`
- `aria-label="Localização na app"` + separador visual `›`

---

## 4. Sistema de temas e paletas

### 4.1 Actual (UI-001)

- `ThemeMode`: `light` | `dark` | `system`
- `data-theme` no `<html>` via `theme.ts`
- Selector em `/settings` e ciclo rápido no topbar

### 4.2 Alvo (UI-013) — modo + paleta

Separar **modo** (claro/escuro/sistema) de **paleta** (identidade visual):

```
Preferência = { mode: ThemeMode, palette: PaletteId }
```

| Paleta | `PaletteId` | Personalidade |
|---|---|---|
| **Aegis** (omissão) | `aegis` | Azul contido, confiança — tokens actuais |
| **Aurora** | `aurora` | Futurista suave — teal/violeta desaturado |
| **Midnight** | `midnight` | Quase monocromático, accent frio — sessões longas |
| **Paper** | `paper` | Claro quente, sombras mínimas — RH/documentos |

Implementação CSS:

```css
/* tokens.css — blocos por paleta */
[data-palette="aurora"][data-theme="dark"] {
  --color-accent: /* … */;
  --color-bg-base: /* … */;
}
```

Persistência: `localStorage` chave `aegis-palette` (paralela a `aegis-theme`).

UI em `/settings` → secção **Aparência**:

- Radiogroup **Modo** (existente)
- Grid de **cartões de paleta** com pré-visualização ao vivo
- O topbar mantém ciclo rápido só de **modo**; paleta só nas definições

> ⚠️ **Segurança:** paletas alteram cor, não layout. Estados de perigo (`--color-danger`)
> mantêm contraste WCAG AA em todas as combinações — testar com ferramenta de contraste
> antes de marcar paleta como disponível.

---

## 5. Biblioteca de componentes (`lib/ui/`)

### 5.1 Problema

O catálogo `/dev/components` (UI-010) documenta tokens mas os botões, forms e
cards são **CSS copiado** na página — não há imports reutilizáveis.

### 5.2 Componentes base (UI-012)

| Componente | Substitui | Props / slots principais |
|---|---|---|
| `PageShell` | `.page`, `.page-head`, `h1`, `.lead` | `title`, `description`, `breadcrumb`, `taskId?`, `actions` slot |
| `Panel` | `.panel`, `.panel-head` | `title`, `padding`, `variant` |
| `Eyebrow` | `.eyebrow` (20+ cópias) | `text`, `devOnly` |
| `Button` | `.btn`, `.primary`, `.secondary`, `.danger` | `variant`, `size`, `loading`, `href?` |
| `Field` | labels + inputs inline | `label`, `error`, `hint`, `type` |
| `EmptyState` | blocos `.empty` | `title`, `description`, `action` slot |
| `HubLinks` | `.links` em `/work`, `/security` | `items: {href, title, description}[]` |
| `StatusBanner` | `.status`, `.error` | `variant: info\|success\|warning\|error` |
| `Section` | `.block` em `/settings` | `title`, `hint` |

### 5.3 Componentes de dados (UI-015)

| Componente | Uso |
|---|---|
| `DataTable` | CRM, fin, admin — modo `dense` para admin |
| `MetricCard` | KPIs fin/crm — `label`, `value`, `trend?` |
| `ListRow` | Cofre, listas densas — título + meta + chevron |

### 5.4 Feedback e acções (UI-017)

| Componente | Uso |
|---|---|
| `Toast` | Copy password, item guardado, erros não bloqueantes |
| `Skeleton` | Carregamento — substituir «A carregar…» |
| `ConfirmDialog` | Wipe, offboarding, acções destrutivas — multi-step |

### 5.5 Convenções Svelte 5

- **Slots** (`{#snippet actions()}`) para acções de página sem props gigantes
- Estilos **scoped** no componente; tokens via `var(--*)` apenas
- Cada componente no catálogo `/dev/components` com exemplo importável
- `prefers-reduced-motion` em todas as transições

### 5.6 Estratégia de migração (UI-016)

Sem *big-bang* — módulo a módulo:

```
1. Extrair de /dev/components → lib/ui/
2. Registar em design/catalog.ts
3. Piloto: /fin/fiscal (página pequena)
4. Módulo: /fin/*
5. Segurança, RH, CRM, Equipa, Admin
6. Cofre (já mais polido — último ou incremental)
```

Critério de «página migrada»: usa `PageShell` + zero CSS local de `.page-head`/`.eyebrow`.

---

## 6. UX, IX e micro-interacções

| Momento | Hoje | Alvo |
|---|---|---|
| Orientação | Eyebrows com task IDs | Breadcrumbs + título; task ID só em DEV |
| Navegação | Sidebar flat + ⌘K flat | Árvore + breadcrumbs + ⌘K hierárquico |
| Teclado | ⌘K global | + atalhos contextuais (`/` pesquisa, `n` novo no cofre) |
| Carregamento | Texto «A carregar…» | `Skeleton` por secção |
| Erros | `<p class="error">` disperso | `StatusBanner` + `Toast` |
| Empty states | Variados | `EmptyState` com CTA consistente |
| Confiança ZK | Copy dispersa | Badge «Desbloqueado localmente» no header quando MK activa |
| Acções destrutivas | Inconsistente | `ConfirmDialog` padrão (wipe, offboarding) |
| Reduced motion | Parcial em shell | Obrigatório em `lib/ui/*` |

Ver também [`product-ui-vision.md` §6](product-ui-vision.md) (timing GSAP, stagger, clipboard).

---

## 7. Auditoria por módulo

| Módulo | Rotas | Prioridade migração | Notas |
|---|---|---|---|
| **Finanças** | `/fin/*` | Alta | Mais CSS duplicado; criar hub `/fin` |
| **RH** | `/hr/*` | Alta | Painéis aninhados → `Panel` + layout 2 colunas |
| **Equipa** | `/team/*` | Média | Página principal longa → hub como `/work` |
| **CRM** | `/crm` | Média | `DataTable` + `MetricCard` |
| **Segurança** | `/security/*` | Média | Hub OK; unificar páginas com CSS solto |
| **Trabalho** | `/work/*` | Baixa | Hub OK; agrupar Google em sub-secção |
| **Cofre** | `/vault/*` | Baixa | Já tem motion; acrescentar breadcrumbs |
| **Admin** | `/admin/*` | Média | Tabelas densas + `ConfirmDialog` wipe |
| **Definições** | `/settings` | Média | Nav lateral + picker de paletas |
| **Auth** | `/auth/*` | Baixa | `AuthShell` isolado; alinhar `Button` |

---

## 8. Tasks de implementação

| ID | Entrega | Depende de | Tamanho |
|---|---|---|---|
| [UI-011](../05-tasks/backlog.md) | `ROUTE_TREE` + sidebar hierárquica + `Breadcrumbs` | UI-002 | L |
| [UI-012](../05-tasks/backlog.md) | `lib/ui/` — PageShell, Panel, Button, Field, EmptyState, HubLinks | UI-001, UI-010 | M |
| [UI-013](../05-tasks/backlog.md) | Paletas múltiplas + picker em `/settings` | UI-001 | M |
| [UI-014](../05-tasks/backlog.md) | Hubs `/fin`, `/team`; limpar sidebar flat | UI-011 | S |
| [UI-015](../05-tasks/backlog.md) | `DataTable`, `MetricCard`, `ListRow` | UI-012 | M |
| [UI-016](../05-tasks/backlog.md) | Migração página-a-página (começar `/fin/*`) | UI-012, UI-011 | L |
| [UI-017](../05-tasks/backlog.md) | `Toast`, `Skeleton`, `ConfirmDialog` | UI-012 | M |

**Ordem recomendada:** UI-012 → UI-011 → UI-013 → UI-014 → UI-015/UI-017 → UI-016.

Pode correr **em paralelo** com produção (VPS, Mail, Google) — não bloqueia infra.

---

## 9. Critérios de aceitação (DoD design system)

- [ ] `lib/ui/` com ≥8 componentes base importáveis
- [ ] Sidebar em árvore; sub-páginas fin não aparecem ao nível de Cofre
- [ ] Breadcrumbs em todas as rotas `(app)/` excepto hubs de primeiro nível
- [ ] ≥4 paletas com contraste WCAG AA validado
- [ ] Catálogo `/dev/components` mostra imports reais, não CSS inline
- [ ] Módulo `/fin/*` 100% migrado para `PageShell`
- [ ] Eyebrows com task ID só visíveis em `import.meta.env.DEV`
- [ ] `prefers-reduced-motion` testado nos componentes `lib/ui/`

---

## 10. Ligações

- North Star: [`product-ui-vision.md`](product-ui-vision.md)
- Backlog: [`../05-tasks/backlog.md`](../05-tasks/backlog.md) (`UI-011`–`UI-017`)
- Código tokens: `frontend/src/lib/design/tokens.css`
- Shell actual: `frontend/src/lib/shell/`
- Catálogo dev: `frontend/src/routes/(app)/dev/components/`
