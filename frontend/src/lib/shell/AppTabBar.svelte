<script lang="ts">
  import { isRouteActive, type RouteNode } from "./routes";
  import { iconForHref } from "./navIcons";
  import NavIcon from "./NavIcon.svelte";

  interface Props {
    items: RouteNode[];
    pathname: string;
  }

  let { items, pathname }: Props = $props();
</script>

<nav class="tab-bar" aria-label="Navegação rápida">
  {#each items as item (item.href)}
    <a
      href={item.href}
      class="tab"
      class:active={isRouteActive(pathname, item.href)}
      aria-current={isRouteActive(pathname, item.href) ? "page" : undefined}
    >
      <span class="tab-icon" aria-hidden="true">
        <NavIcon name={iconForHref(item.href)} size={18} />
      </span>
      <span class="tab-label">{item.label}</span>
    </a>
  {/each}
</nav>

<style>
  .tab-bar {
    display: flex;
    position: fixed;
    bottom: 0;
    left: 0;
    right: 0;
    z-index: 50;
    height: var(--shell-tab-height, 3.5rem);
    padding: var(--space-1) var(--space-2);
    padding-bottom: max(var(--space-1), env(safe-area-inset-bottom));
    border-top: 1px solid var(--color-border);
    background: var(--color-bg-elevated);
    box-sizing: border-box;
  }

  @media (min-width: 768px) {
    .tab-bar {
      display: none;
    }
  }

  .tab {
    flex: 1;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    gap: 2px;
    min-height: 44px;
    border-radius: var(--radius-sm);
    color: var(--color-text-muted);
    text-decoration: none;
    font-size: var(--text-xs);
    transition: color var(--duration-fast) var(--ease-out);
  }

  .tab.active {
    color: var(--color-accent);
  }

  .tab-icon {
    font-family: var(--font-display);
    font-weight: 600;
    font-size: var(--text-sm);
    line-height: 1;
  }

  .tab-label {
    line-height: 1.1;
  }

  @media (prefers-reduced-motion: reduce) {
    .tab {
      transition: none;
    }
  }
</style>
