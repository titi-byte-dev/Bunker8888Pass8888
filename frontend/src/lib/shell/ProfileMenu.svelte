<script lang="ts">
  /**
   * ProfileMenu — avatar no topbar com dropdown (Definições, Sair).
   * Didático: Definições saiu da sidebar; acções de conta ficam aqui.
   */
  import { onMount } from "svelte";
  import { page } from "$app/state";

  interface Props {
    email?: string | null;
    onLogout: () => void;
  }

  let { email = null, onLogout }: Props = $props();

  let open = $state(false);
  let rootEl: HTMLDivElement | undefined = $state();

  const initial = $derived(
    email ? email.trim().charAt(0).toUpperCase() || "?" : "?",
  );
  const settingsActive = $derived(page.url.pathname.startsWith("/settings"));

  onMount(() => {
    function onDocumentClick(e: MouseEvent) {
      if (!open || !rootEl) return;
      if (!rootEl.contains(e.target as Node)) open = false;
    }
    function onKey(e: KeyboardEvent) {
      if (e.key === "Escape") open = false;
    }
    document.addEventListener("click", onDocumentClick);
    document.addEventListener("keydown", onKey);
    return () => {
      document.removeEventListener("click", onDocumentClick);
      document.removeEventListener("keydown", onKey);
    };
  });

  function toggle(e: MouseEvent) {
    e.stopPropagation();
    open = !open;
  }

  function close() {
    open = false;
  }

  function logout() {
    close();
    onLogout();
  }
</script>

<div class="profile-menu" bind:this={rootEl}>
  <button
    type="button"
    class="profile-trigger"
    class:active={open || settingsActive}
    onclick={toggle}
    aria-expanded={open}
    aria-haspopup="menu"
    aria-label="Menu de conta"
    title="Conta e definições"
  >
    <span class="avatar" aria-hidden="true">{initial}</span>
  </button>

  {#if open}
    <div class="dropdown" role="menu">
      {#if email}
        <p class="dropdown-email" title={email}>{email}</p>
      {/if}

      <a href="/settings" class="dropdown-item" role="menuitem" onclick={close}>
        <svg width="18" height="18" viewBox="0 0 24 24" aria-hidden="true" focusable="false">
          <path
            fill="currentColor"
            d="M12 8a4 4 0 1 1 0 8 4 4 0 0 1 0-8zm8.94 3a7.96 7.96 0 0 0 .05-.94 7.96 7.96 0 0 0-.05-.94l2.03-1.58a.5.5 0 0 0 .12-.64l-1.92-3.32a.5.5 0 0 0-.6-.22l-2.39.96a7.28 7.28 0 0 0-1.62-.94l-.36-2.54A.5.5 0 0 0 14 2h-4a.5.5 0 0 0-.49.42l-.36 2.54c-.58.22-1.12.52-1.62.94l-2.39-.96a.5.5 0 0 0-.6.22L2.62 8.9a.5.5 0 0 0 .12.64L4.77 11.1c-.03.31-.05.63-.05.94s.02.63.05.94l-2.03 1.58a.5.5 0 0 0-.12.64l1.92 3.32a.5.5 0 0 0 .6.22l2.39-.96c.5.42 1.04.77 1.62.99l.36 2.54A.5.5 0 0 0 10 22h4a.5.5 0 0 0 .49-.42l.36-2.54a7.28 7.28 0 0 0 1.62-.99l2.39.96a.5.5 0 0 0 .6-.22l1.92-3.32a.5.5 0 0 0-.12-.64l-2.03-1.58z"
          />
        </svg>
        <span>Definições</span>
      </a>

      <div class="dropdown-divider" role="separator"></div>

      <button type="button" class="dropdown-item danger" role="menuitem" onclick={logout}>
        <svg width="18" height="18" viewBox="0 0 24 24" aria-hidden="true" focusable="false">
          <path
            fill="currentColor"
            d="M17 7l-1.41 1.41L18.17 11H8v2h10.17l-2.58 2.58L17 17l5-5-5-5zM4 5h8V3H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h8v-2H4V5z"
          />
        </svg>
        <span>Sair</span>
      </button>
    </div>
  {/if}
</div>

<style>
  .profile-menu {
    position: relative;
  }

  .profile-trigger {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 2.25rem;
    height: 2.25rem;
    padding: 0;
    border: 1px solid var(--color-border);
    border-radius: var(--radius-full, 9999px);
    background: var(--color-bg-surface);
    color: var(--color-text);
    cursor: pointer;
    transition:
      border-color var(--duration-fast) var(--ease-out),
      background-color var(--duration-fast) var(--ease-out),
      box-shadow var(--duration-fast) var(--ease-out);
  }

  .profile-trigger:hover,
  .profile-trigger.active {
    border-color: var(--color-accent);
    background: var(--color-accent-muted);
    color: var(--color-accent);
  }

  .profile-trigger:focus-visible {
    outline: 2px solid var(--color-accent);
    outline-offset: 2px;
  }

  .avatar {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 100%;
    height: 100%;
    font-family: var(--font-display);
    font-size: var(--text-sm);
    font-weight: 600;
    line-height: 1;
  }

  .dropdown {
    position: absolute;
    top: calc(100% + var(--space-2));
    right: 0;
    z-index: 60;
    min-width: 13rem;
    padding: var(--space-2);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-md, var(--radius-sm));
    background: var(--color-bg-elevated);
    box-shadow: 0 8px 24px color-mix(in srgb, var(--color-bg-base) 35%, transparent);
  }

  .dropdown-email {
    margin: 0 0 var(--space-1);
    padding: var(--nav-item-py) var(--nav-item-px);
    font-size: var(--text-xs);
    line-height: var(--nav-item-leading);
    color: var(--color-text-muted);
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  .dropdown-item {
    display: flex;
    align-items: center;
    gap: var(--space-2);
    width: 100%;
    padding: var(--nav-item-py) var(--nav-item-px);
    border: none;
    border-radius: var(--radius-sm);
    background: transparent;
    color: var(--color-text);
    font-family: var(--font-ui);
    font-size: var(--text-sm);
    line-height: var(--nav-item-leading);
    min-height: var(--nav-item-min-height);
    text-decoration: none;
    text-align: left;
    cursor: pointer;
    box-sizing: border-box;
  }

  .dropdown-item:hover {
    background: var(--color-bg-surface);
  }

  .dropdown-item.danger {
    color: var(--color-danger, #b42318);
  }

  .dropdown-item.danger:hover {
    background: color-mix(in srgb, var(--color-danger, #b42318) 8%, transparent);
  }

  .dropdown-divider {
    height: 1px;
    margin: var(--space-1) 0;
    background: var(--color-border);
  }

  @media (prefers-reduced-motion: reduce) {
    .profile-trigger {
      transition: none;
    }
  }
</style>
