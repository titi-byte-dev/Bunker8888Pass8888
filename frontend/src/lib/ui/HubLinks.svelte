<script lang="ts">
  /**
   * HubLinks (UI-012) — grelha de cartoes-link para paginas-hub
   * (/work, /security, /fin). Substitui `.links` repetido.
   */
  export type HubLinkItem = {
    href: string;
    title: string;
    description?: string;
    /** ID de task — so em DEV. */
    taskId?: string;
    comingSoon?: boolean;
  };

  interface Props {
    items: HubLinkItem[];
  }
  let { items }: Props = $props();
  const dev = import.meta.env.DEV;
</script>

<ul class="hub">
  {#each items as item (item.href)}
    <li>
      <svelte:element
        this={item.comingSoon ? "span" : "a"}
        class="card"
        class:disabled={item.comingSoon}
        href={item.comingSoon ? undefined : item.href}
        aria-disabled={item.comingSoon ? "true" : undefined}
      >
        <span class="row">
          <span class="title">{item.title}</span>
          {#if item.comingSoon}<span class="badge">Em breve</span>{/if}
          {#if dev && item.taskId}<span class="task">{item.taskId}</span>{/if}
        </span>
        {#if item.description}<span class="desc">{item.description}</span>{/if}
      </svelte:element>
    </li>
  {/each}
</ul>

<style>
  .hub {
    list-style: none;
    margin: 0;
    padding: 0;
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(15rem, 1fr));
    gap: var(--space-2);
  }
  .card {
    display: flex;
    flex-direction: column;
    gap: 0.125rem;
    padding: var(--space-3);
    height: auto;
    box-sizing: border-box;
    background: var(--color-bg-surface);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-md);
    text-decoration: none;
    color: var(--color-text);
    transition:
      border-color var(--duration-fast) var(--ease-out),
      transform var(--duration-fast) var(--ease-out);
  }
  a.card:hover {
    border-color: var(--color-border-strong);
    transform: translateY(-1px);
  }
  .card.disabled { opacity: 0.55; cursor: not-allowed; }
  .row { display: flex; align-items: center; gap: var(--space-2); }
  .title { font-weight: 600; font-size: var(--text-sm); line-height: var(--nav-item-leading); }
  .desc { color: var(--color-text-muted); font-size: var(--text-xs); line-height: var(--leading-snug); }
  .badge, .task {
    font-size: var(--text-xs);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-sm);
    padding: 0 var(--space-1);
    color: var(--color-text-muted);
  }
  .task { font-family: var(--font-mono); margin-left: auto; }
  @media (prefers-reduced-motion: reduce) {
    .card { transition: none; }
    a.card:hover { transform: none; }
  }
</style>
