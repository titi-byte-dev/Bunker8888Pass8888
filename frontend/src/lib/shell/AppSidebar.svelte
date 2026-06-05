<script lang="ts">
  /**
   * AppSidebar (UI-011) — navegacao em ARVORE.
   * Antes: lista plana (~15 links irmaos). Agora: modulos de topo + filhos
   * que se expandem quando o modulo (ou um filho) esta activo.
   * Fonte unica: lib/shell/routes.ts (ROUTE_TREE).
   */
  import { isRouteActive, type RouteNode } from "./routes";

  interface Props {
    modules: RouteNode[];
    pathname: string;
  }
  let { modules, pathname }: Props = $props();

  const dev = import.meta.env.DEV;

  /** Um modulo expande-se se ele ou algum filho corresponder ao pathname. */
  function moduleOpen(node: RouteNode): boolean {
    if (isRouteActive(pathname, node.href)) return true;
    return (node.children ?? []).some((c) => isRouteActive(pathname, c.href));
  }
</script>

<nav class="sidebar" aria-label="Navegacao principal">
  <div class="brand">
    <span class="brand-mark" aria-hidden="true">◆</span>
    <span class="brand-name">AegisPass</span>
  </div>

  <ul class="nav-list">
    {#each modules as mod (mod.href)}
      {@const open = moduleOpen(mod)}
      <li>
        {#if mod.comingSoon}
          <span class="nav-link disabled" aria-disabled="true">
            {mod.label}<span class="badge">Em breve</span>
          </span>
        {:else}
          <a
            href={mod.href}
            class="nav-link module"
            class:active={isRouteActive(pathname, mod.href)}
            aria-current={pathname === mod.href ? "page" : undefined}
          >
            {mod.label}
            {#if dev && mod.taskId}<span class="task">{mod.taskId}</span>{/if}
          </a>
        {/if}

        {#if mod.children && open}
          <ul class="children">
            {#each mod.children as child (child.href)}
              <a
                href={child.href}
                class="nav-link child"
                class:active={pathname === child.href}
                aria-current={pathname === child.href ? "page" : undefined}
              >
                {child.label}
              </a>
            {/each}
          </ul>
        {/if}
      </li>
    {/each}
  </ul>
</nav>

<style>
  .sidebar {
    display: none;
    flex-direction: column;
    width: var(--shell-sidebar-width, 15rem);
    min-height: 100dvh;
    padding: var(--space-4);
    border-right: 1px solid var(--color-border);
    background: var(--color-bg-elevated);
    box-sizing: border-box;
  }
  @media (min-width: 768px) {
    .sidebar { display: flex; }
  }

  .brand {
    display: flex;
    align-items: center;
    gap: var(--space-2);
    padding: var(--space-2) var(--space-2) var(--space-6);
  }
  .brand-mark { color: var(--color-accent); font-size: var(--text-lg); }
  .brand-name {
    font-family: var(--font-display);
    font-weight: 600;
    font-size: var(--text-lg);
    line-height: var(--leading-tight);
  }

  .nav-list {
    list-style: none;
    margin: 0;
    padding: 0;
    display: flex;
    flex-direction: column;
    gap: var(--space-1);
  }

  .nav-link {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: var(--space-2);
    padding: var(--space-2) var(--space-3);
    border-radius: var(--radius-sm);
    color: var(--color-text-muted);
    text-decoration: none;
    font-size: var(--text-sm);
    font-weight: 500;
    transition:
      background-color var(--duration-fast) var(--ease-out),
      color var(--duration-fast) var(--ease-out);
  }
  .nav-link.module { font-weight: 600; }
  .nav-link:hover:not(.disabled) {
    background: var(--color-bg-surface);
    color: var(--color-text);
  }
  .nav-link.active { color: var(--color-accent); }
  .nav-link.module.active { background: var(--color-accent-muted); }
  .nav-link.disabled { opacity: 0.55; cursor: not-allowed; }

  .children {
    list-style: none;
    margin: var(--space-1) 0 var(--space-1) var(--space-3);
    padding-left: var(--space-3);
    border-left: 1px solid var(--color-border);
    display: flex;
    flex-direction: column;
    gap: 2px;
  }
  .nav-link.child {
    font-weight: 500;
    font-size: var(--text-xs);
    padding: var(--space-1) var(--space-2);
  }
  .nav-link.child.active {
    background: var(--color-accent-muted);
    color: var(--color-accent);
  }

  .badge, .task {
    font-size: var(--text-xs);
    font-weight: 500;
    color: var(--color-text-muted);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-sm);
    padding: 0 var(--space-1);
  }
  .task { font-family: var(--font-mono); }

  @media (prefers-reduced-motion: reduce) {
    .nav-link { transition: none; }
  }
</style>
