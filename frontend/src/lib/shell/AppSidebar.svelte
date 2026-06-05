<script lang="ts">
  /**
   * AppSidebar (UI-011) — navegação em árvore com ícones e modo colapsável.
   * Colapsada: só ícones (~3.75rem) para libertar espaço ao conteúdo em desktop.
   */
  import { isRouteActive, type RouteNode } from "./routes";
  import { iconForHref } from "./navIcons";
  import NavIcon from "./NavIcon.svelte";

  interface Props {
    modules: RouteNode[];
    pathname: string;
    collapsed?: boolean;
  }

  let { modules, pathname, collapsed = false }: Props = $props();

  const dev = import.meta.env.DEV;

  function moduleOpen(node: RouteNode): boolean {
    if (collapsed) return false;
    if (isRouteActive(pathname, node.href)) return true;
    return (node.children ?? []).some((c) => isRouteActive(pathname, c.href));
  }
</script>

<nav class="sidebar" class:collapsed aria-label="Navegação principal">
  <div class="brand">
    <span class="brand-mark" aria-hidden="true">◆</span>
    {#if !collapsed}
      <span class="brand-name">AegisPass</span>
    {/if}
  </div>

  <ul class="nav-list">
    {#each modules as mod (mod.href)}
      {@const open = moduleOpen(mod)}
      <li>
        {#if mod.comingSoon}
          <span class="nav-link disabled" aria-disabled="true" title={mod.label}>
            <NavIcon name={iconForHref(mod.href)} />
            {#if !collapsed}
              <span class="nav-text">{mod.label}</span>
              <span class="badge">Em breve</span>
            {/if}
          </span>
        {:else}
          <a
            href={mod.href}
            class="nav-link module"
            class:active={isRouteActive(pathname, mod.href)}
            aria-current={pathname === mod.href ? "page" : undefined}
            title={collapsed ? mod.label : undefined}
          >
            <NavIcon name={iconForHref(mod.href)} />
            {#if !collapsed}
              <span class="nav-text">{mod.label}</span>
              {#if dev && mod.taskId}<span class="task">{mod.taskId}</span>{/if}
            {/if}
          </a>
        {/if}

        {#if mod.children && open}
          <ul class="children">
            {#each mod.children as child (child.href)}
              <li>
                <a
                  href={child.href}
                  class="nav-link child"
                  class:active={pathname === child.href}
                  aria-current={pathname === child.href ? "page" : undefined}
                >
                  <NavIcon name="child" size={14} />
                  <span class="nav-text">{child.label}</span>
                </a>
              </li>
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
    transition: width var(--duration-normal) var(--ease-out);
    flex-shrink: 0;
  }

  .sidebar.collapsed {
    width: var(--shell-sidebar-collapsed-width, 3.75rem);
    padding-inline: var(--space-2);
    align-items: center;
  }

  @media (min-width: 768px) {
    .sidebar {
      display: flex;
    }
  }

  .brand {
    display: flex;
    align-items: center;
    gap: var(--space-2);
    padding: var(--space-2) var(--space-2) var(--space-6);
    width: 100%;
  }

  .sidebar.collapsed .brand {
    justify-content: center;
    padding-bottom: var(--space-4);
  }

  .brand-mark {
    color: var(--color-accent);
    font-size: var(--text-lg);
    flex-shrink: 0;
  }

  .brand-name {
    font-family: var(--font-display);
    font-weight: 600;
    font-size: var(--text-lg);
    line-height: var(--leading-tight);
    white-space: nowrap;
  }

  .nav-list {
    list-style: none;
    margin: 0;
    padding: 0;
    display: flex;
    flex-direction: column;
    gap: var(--space-1);
    width: 100%;
    flex: 1;
    min-height: 0;
    overflow-y: auto;
  }

  .nav-link {
    display: flex;
    align-items: center;
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
    min-height: 2.25rem;
  }

  .sidebar.collapsed .nav-link {
    justify-content: center;
    padding-inline: var(--space-2);
  }

  .nav-link.module {
    font-weight: 600;
  }

  .nav-text {
    flex: 1;
    min-width: 0;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .nav-link:hover:not(.disabled) {
    background: var(--color-bg-surface);
    color: var(--color-text);
  }

  .nav-link.active {
    color: var(--color-accent);
  }

  .nav-link.module.active {
    background: var(--color-accent-muted);
  }

  .nav-link.disabled {
    opacity: 0.55;
    cursor: not-allowed;
  }

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

  .badge,
  .task {
    font-size: var(--text-xs);
    font-weight: 500;
    color: var(--color-text-muted);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-sm);
    padding: 0 var(--space-1);
    flex-shrink: 0;
  }

  .task {
    font-family: var(--font-mono);
  }

  @media (prefers-reduced-motion: reduce) {
    .sidebar,
    .nav-link {
      transition: none;
    }
  }
</style>
