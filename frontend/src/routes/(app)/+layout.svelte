<script lang="ts">
  import { onMount } from "svelte";
  import { goto } from "$app/navigation";
  import { page } from "$app/state";
  import AppSidebar from "$lib/shell/AppSidebar.svelte";
  import AppTabBar from "$lib/shell/AppTabBar.svelte";
  import CommandPalette from "$lib/shell/CommandPalette.svelte";
  import { navModules, tabBarModules } from "$lib/shell/routes";
  import {
    loadThemePreference,
    setThemePreference,
    cycleThemePreference,
    themeModeLabel,
    type ThemeMode,
  } from "$lib/design";
  import { clearSession, loadUserEmail } from "$lib/session";
  import { purgeMasterKey } from "$lib/vault/masterKeyStore";
  import ShellPageMotion from "$lib/motion/ShellPageMotion.svelte";
  import { ToastHost, ConfirmDialog } from "$lib/ui";
  import { loadSidebarCollapsed, saveSidebarCollapsed } from "$lib/shell/sidebarState";

  let { children } = $props();

  let themeMode = $state<ThemeMode>(loadThemePreference());
  let paletteOpen = $state(false);
  let sidebarCollapsed = $state(loadSidebarCollapsed());

  onMount(() => {
    function onKey(e: KeyboardEvent) {
      if ((e.metaKey || e.ctrlKey) && e.key.toLowerCase() === "k") {
        e.preventDefault();
        paletteOpen = !paletteOpen;
      }
    }
    window.addEventListener("keydown", onKey);
    return () => window.removeEventListener("keydown", onKey);
  });

  function toggleTheme() {
    themeMode = cycleThemePreference(themeMode);
    setThemePreference(themeMode);
  }

  function handleLogout() {
    purgeMasterKey();
    clearSession();
    goto("/auth/login");
  }

  const navItems = navModules();
  const tabs = tabBarModules();
  const userEmail = $derived(loadUserEmail());

  function toggleSidebar() {
    sidebarCollapsed = !sidebarCollapsed;
    saveSidebarCollapsed(sidebarCollapsed);
  }
</script>

<div class="app-shell" class:sidebar-collapsed={sidebarCollapsed}>
  <AppSidebar modules={navItems} pathname={page.url.pathname} collapsed={sidebarCollapsed} />

  <!-- Separador com botão no topo — entre sidebar e conteúdo (só desktop). -->
  <div class="shell-divider">
    <button
      type="button"
      class="sidebar-toggle"
      onclick={toggleSidebar}
      aria-expanded={!sidebarCollapsed}
      aria-label={sidebarCollapsed ? "Expandir menu lateral" : "Recolher menu lateral"}
      title={sidebarCollapsed ? "Expandir menu" : "Recolher menu"}
    >
      <span aria-hidden="true">{sidebarCollapsed ? "›" : "‹"}</span>
    </button>
  </div>

  <div class="shell-main">
    <header class="topbar">
      {#if userEmail}
        <span class="user-email" title="Sessão activa">{userEmail}</span>
      {:else}
        <div class="topbar-spacer" aria-hidden="true"></div>
      {/if}
      <div class="topbar-actions">
        <button type="button" class="palette-btn" onclick={() => (paletteOpen = true)} title="Command palette (Ctrl+K)">
          ⌘K
        </button>
        <button type="button" class="theme-btn" onclick={toggleTheme} title="Alternar tema">
          {themeModeLabel(themeMode)}
        </button>
        <button type="button" class="logout-btn" onclick={handleLogout} title="Terminar sessão">
          Sair
        </button>
      </div>
    </header>

    <main class="shell-content">
      <ShellPageMotion>
        {@render children()}
      </ShellPageMotion>
    </main>
  </div>

  <AppTabBar items={tabs} pathname={page.url.pathname} />
</div>

<CommandPalette open={paletteOpen} onClose={() => (paletteOpen = false)} />
<ToastHost />
<ConfirmDialog />

<style>
  .app-shell {
    --shell-sidebar-width: 15rem;
    --shell-sidebar-collapsed-width: 3.75rem;
    --shell-tab-height: 3.5rem;
    --shell-toggle-top: 3.25rem;
    display: flex;
    min-height: 100dvh;
    position: relative;
  }

  /* Linha zero entre sidebar e main — o botão centra-se na fronteira. */
  .shell-divider {
    display: none;
    flex-shrink: 0;
    width: 0;
    position: relative;
    z-index: 30;
  }

  @media (min-width: 768px) {
    .shell-divider {
      display: block;
    }
  }

  .sidebar-toggle {
    position: absolute;
    left: 0;
    top: var(--shell-toggle-top);
    transform: translateX(-50%);
    display: flex;
    align-items: center;
    justify-content: center;
    width: 1.625rem;
    height: 2.75rem;
    padding: 0;
    border: 1px solid var(--color-border);
    border-radius: var(--radius-sm);
    background: var(--color-bg-elevated);
    color: var(--color-text-muted);
    font-size: var(--text-lg);
    line-height: 1;
    cursor: pointer;
    box-shadow: 0 1px 4px color-mix(in srgb, var(--color-bg-base) 40%, transparent);
    transition:
      background-color var(--duration-fast) var(--ease-out),
      color var(--duration-fast) var(--ease-out),
      border-color var(--duration-fast) var(--ease-out);
  }

  .sidebar-toggle:hover {
    color: var(--color-accent);
    border-color: var(--color-accent);
    background: var(--color-accent-muted);
  }

  .sidebar-toggle:focus-visible {
    outline: 2px solid var(--color-accent);
    outline-offset: 2px;
  }

  .shell-main {
    flex: 1;
    display: flex;
    flex-direction: column;
    min-width: 0;
  }

  .topbar {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: var(--space-3);
    padding: var(--space-3) var(--space-4);
    border-bottom: 1px solid var(--color-border);
    background: var(--color-bg-base);
  }

  .user-email {
    font-size: var(--text-sm);
    color: var(--color-text-muted);
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    min-width: 0;
  }

  .topbar-actions {
    display: flex;
    align-items: center;
    gap: var(--space-2);
    flex-shrink: 0;
  }

  .logout-btn {
    padding: var(--space-2) var(--space-3);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-sm);
    background: transparent;
    color: var(--color-text-muted);
    font-family: var(--font-ui);
    font-size: var(--text-sm);
    cursor: pointer;
  }

  .logout-btn:hover {
    color: var(--color-text);
    background: var(--color-bg-surface);
  }

  .palette-btn {
    padding: var(--space-2) var(--space-3);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-sm);
    background: transparent;
    color: var(--color-text-muted);
    font-family: var(--font-mono);
    font-size: var(--text-xs);
    cursor: pointer;
  }

  .palette-btn:hover {
    color: var(--color-text);
    background: var(--color-bg-surface);
  }

  .theme-btn {
    padding: var(--space-2) var(--space-3);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-sm);
    background: var(--color-bg-surface);
    color: var(--color-text);
    font-family: var(--font-ui);
    font-size: var(--text-sm);
    cursor: pointer;
    transition: background-color var(--duration-fast) var(--ease-out);
  }

  .theme-btn:hover {
    background: var(--color-accent-muted);
  }

  .shell-content {
    flex: 1;
    padding: var(--space-6) var(--space-4);
    padding-bottom: calc(var(--shell-tab-height) + var(--space-8));
    width: 100%;
    max-width: none;
    margin: 0;
    box-sizing: border-box;
  }

  @media (min-width: 768px) {
    .shell-content {
      padding-bottom: var(--space-12);
      padding-inline: var(--space-6);
    }
  }

  @media (prefers-reduced-motion: reduce) {
    .theme-btn,
    .sidebar-toggle {
      transition: none;
    }
  }
</style>
