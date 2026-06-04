<script lang="ts">
  import { page } from "$app/state";
  import AppSidebar from "$lib/shell/AppSidebar.svelte";
  import AppTabBar from "$lib/shell/AppTabBar.svelte";
  import { tabBarItems, visibleNavItems } from "$lib/shell/nav";
  import {
    loadThemePreference,
    setThemePreference,
    cycleThemePreference,
    themeModeLabel,
    type ThemeMode,
  } from "$lib/design";

  let { children } = $props();

  let themeMode = $state<ThemeMode>(loadThemePreference());

  function toggleTheme() {
    themeMode = cycleThemePreference(themeMode);
    setThemePreference(themeMode);
  }

  const navItems = visibleNavItems();
  const tabs = tabBarItems();
</script>

<div class="app-shell">
  <AppSidebar items={navItems} pathname={page.url.pathname} />

  <div class="shell-main">
    <header class="topbar">
      <div class="topbar-spacer" aria-hidden="true"></div>
      <button type="button" class="theme-btn" onclick={toggleTheme} title="Alternar tema">
        {themeModeLabel(themeMode)}
      </button>
    </header>

    <main class="shell-content">
      {@render children()}
    </main>
  </div>

  <AppTabBar items={tabs} pathname={page.url.pathname} />
</div>

<style>
  .app-shell {
    --shell-sidebar-width: 15rem;
    --shell-tab-height: 3.5rem;
    display: flex;
    min-height: 100dvh;
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
    justify-content: flex-end;
    padding: var(--space-3) var(--space-4);
    border-bottom: 1px solid var(--color-border);
    background: var(--color-bg-base);
  }

  @media (min-width: 768px) {
    .topbar-spacer {
      flex: 1;
    }
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
    max-width: 56rem;
    width: 100%;
    margin: 0 auto;
    box-sizing: border-box;
  }

  @media (min-width: 768px) {
    .shell-content {
      padding-bottom: var(--space-12);
    }
  }

  @media (prefers-reduced-motion: reduce) {
    .theme-btn {
      transition: none;
    }
  }
</style>
