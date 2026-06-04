<script lang="ts">
  import type { NavItem } from "./nav";
  import { isNavActive } from "./nav";

  interface Props {
    items: NavItem[];
    pathname: string;
  }

  let { items, pathname }: Props = $props();
</script>

<nav class="sidebar" aria-label="Navegação principal">
  <div class="brand">
    <span class="brand-mark" aria-hidden="true">◆</span>
    <span class="brand-name">AegisPass</span>
  </div>

  <ul class="nav-list">
    {#each items as item (item.id)}
      <li>
        {#if item.comingSoon}
          <span class="nav-link disabled" aria-disabled="true">
            {item.label}
            <span class="badge">Em breve</span>
          </span>
        {:else}
          <a
            href={item.href}
            class="nav-link"
            class:active={isNavActive(pathname, item.href)}
            aria-current={isNavActive(pathname, item.href) ? "page" : undefined}
          >
            {item.label}
          </a>
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
    .sidebar {
      display: flex;
    }
  }

  .brand {
    display: flex;
    align-items: center;
    gap: var(--space-2);
    padding: var(--space-2) var(--space-2) var(--space-6);
  }

  .brand-mark {
    color: var(--color-accent);
    font-size: var(--text-lg);
  }

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

  .nav-link:hover:not(.disabled) {
    background: var(--color-bg-surface);
    color: var(--color-text);
  }

  .nav-link.active {
    background: var(--color-accent-muted);
    color: var(--color-accent);
  }

  .nav-link.disabled {
    opacity: 0.55;
    cursor: not-allowed;
  }

  .badge {
    font-size: var(--text-xs);
    font-weight: 500;
    color: var(--color-text-muted);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-sm);
    padding: 0 var(--space-1);
  }

  @media (prefers-reduced-motion: reduce) {
    .nav-link {
      transition: none;
    }
  }
</style>
