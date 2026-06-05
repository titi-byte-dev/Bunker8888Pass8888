<script lang="ts">
  /**
   * Command palette global (UI-006) — ⌘K / Ctrl+K.
   */
  import {
    buildActionCommands,
    buildNavigationCommands,
    buildVaultCommands,
    filterCommands,
    groupCommands,
    groupLabel,
    type CommandEntry,
  } from "./commands";
  import { buildDocSearchCommands } from "$lib/docs/search";
  import { goto } from "$app/navigation";
  import { loadDecodedLogins } from "$lib/vault/ui";

  interface Props {
    open?: boolean;
    onClose?: () => void;
  }

  let { open = false, onClose }: Props = $props();

  let query = $state("");
  let activeIndex = $state(0);
  let vaultCommands = $state<CommandEntry[]>([]);
  let loadingVault = $state(false);
  let inputEl = $state<HTMLInputElement | undefined>(undefined);

  const docCommands = $derived(
    query.trim() ? buildDocSearchCommands(query) : [],
  );

  const allCommands = $derived([
    ...docCommands,
    ...buildNavigationCommands(),
    ...vaultCommands,
    ...filterCommands(buildActionCommands(), query),
  ]);

  const filtered = $derived(
    query.trim()
      ? [
          ...docCommands,
          ...filterCommands(
            [...buildNavigationCommands(), ...vaultCommands, ...buildActionCommands()],
            query,
          ),
        ]
      : allCommands,
  );
  const grouped = $derived(groupCommands(filtered));
  const flatFiltered = $derived(filtered);

  $effect(() => {
    if (open) {
      query = "";
      activeIndex = 0;
      loadVault();
      queueMicrotask(() => inputEl?.focus());
    }
  });

  $effect(() => {
    if (activeIndex >= flatFiltered.length) {
      activeIndex = Math.max(0, flatFiltered.length - 1);
    }
  });

  async function loadVault() {
    loadingVault = true;
    try {
      const logins = await loadDecodedLogins();
      vaultCommands = buildVaultCommands(logins);
    } catch {
      vaultCommands = [];
    } finally {
      loadingVault = false;
    }
  }

  function close() {
    onClose?.();
  }

  async function runCommand(cmd: CommandEntry) {
    close();
    await goto(cmd.href);
  }

  function onKeydown(e: KeyboardEvent) {
    if (e.key === "Escape") {
      e.preventDefault();
      close();
      return;
    }
    if (e.key === "ArrowDown") {
      e.preventDefault();
      activeIndex = Math.min(activeIndex + 1, flatFiltered.length - 1);
      return;
    }
    if (e.key === "ArrowUp") {
      e.preventDefault();
      activeIndex = Math.max(activeIndex - 1, 0);
      return;
    }
    if (e.key === "Enter" && flatFiltered[activeIndex]) {
      e.preventDefault();
      void runCommand(flatFiltered[activeIndex]!);
    }
  }

  function globalIndex(cmd: CommandEntry): number {
    return flatFiltered.indexOf(cmd);
  }
</script>

{#if open}
  <!-- svelte-ignore a11y_no_noninteractive_element_interactions -->
  <div
    class="palette-backdrop"
    role="presentation"
    onclick={close}
    onkeydown={(e) => e.key === "Escape" && close()}
  >
    <!-- svelte-ignore a11y_click_events_have_key_events a11y_no_static_element_interactions -->
    <div class="palette" role="dialog" aria-modal="true" aria-label="Command palette" tabindex="-1" onclick={(e) => e.stopPropagation()}>
      <div class="search-row">
        <span class="icon" aria-hidden="true">⌘</span>
        <input
          bind:this={inputEl}
          type="search"
          placeholder="Pesquisar docs, páginas, logins, acções…"
          bind:value={query}
          onkeydown={onKeydown}
          autocomplete="off"
          spellcheck="false"
        />
        <kbd class="hint">esc</kbd>
      </div>

      {#if loadingVault}
        <p class="status">A carregar cofre…</p>
      {/if}

      <div class="results" role="listbox">
        {#if flatFiltered.length === 0}
          <p class="empty">Sem resultados.</p>
        {:else}
          {#each [...grouped.entries()] as [group, cmds] (group)}
            <p class="group-label">{groupLabel(group)}</p>
            <ul>
              {#each cmds as cmd (cmd.id)}
                {@const idx = globalIndex(cmd)}
                <li>
                  <button
                    type="button"
                    class:active={idx === activeIndex}
                    role="option"
                    aria-selected={idx === activeIndex}
                    onclick={() => runCommand(cmd)}
                    onmouseenter={() => (activeIndex = idx)}
                  >
                    <span class="label">{cmd.label}</span>
                    <span class="path">{cmd.href}</span>
                  </button>
                </li>
              {/each}
            </ul>
          {/each}
        {/if}
      </div>

      <footer class="footer">
        <span><kbd>↑↓</kbd> navegar</span>
        <span><kbd>↵</kbd> abrir</span>
      </footer>
    </div>
  </div>
{/if}

<style>
  .palette-backdrop {
    position: fixed;
    inset: 0;
    z-index: 200;
    display: flex;
    align-items: flex-start;
    justify-content: center;
    padding: min(12vh, 6rem) var(--space-4) var(--space-4);
    background: color-mix(in srgb, var(--color-bg-base) 55%, transparent);
    backdrop-filter: blur(4px);
  }

  .palette {
    width: min(32rem, 100%);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-lg);
    background: var(--color-bg-elevated);
    box-shadow: 0 16px 48px color-mix(in srgb, #000 35%, transparent);
    overflow: hidden;
  }

  .search-row {
    display: flex;
    align-items: center;
    gap: var(--space-2);
    padding: var(--space-3) var(--space-4);
    border-bottom: 1px solid var(--color-border);
  }

  .icon {
    color: var(--color-text-muted);
    font-size: var(--text-sm);
  }

  .search-row input {
    flex: 1;
    border: none;
    background: transparent;
    color: inherit;
    font-size: var(--text-base);
    font-family: var(--font-ui);
    outline: none;
  }

  .hint {
    font-size: var(--text-xs);
    padding: 2px 6px;
    border-radius: 4px;
    border: 1px solid var(--color-border);
    color: var(--color-text-muted);
    font-family: var(--font-mono);
  }

  .status,
  .empty {
    margin: 0;
    padding: var(--space-4);
    color: var(--color-text-muted);
    font-size: var(--text-sm);
  }

  .results {
    max-height: min(50vh, 24rem);
    overflow-y: auto;
    padding: var(--space-2) 0;
  }

  .group-label {
    margin: var(--space-2) var(--space-4) var(--space-1);
    font-size: var(--text-xs);
    text-transform: uppercase;
    letter-spacing: 0.04em;
    color: var(--color-text-muted);
  }

  ul {
    list-style: none;
    margin: 0;
    padding: 0;
  }

  button {
    width: 100%;
    display: flex;
    justify-content: space-between;
    align-items: center;
    gap: var(--space-3);
    padding: var(--space-2) var(--space-4);
    border: none;
    background: transparent;
    color: inherit;
    font-family: var(--font-ui);
    font-size: var(--text-sm);
    text-align: left;
    cursor: pointer;
  }

  button:hover,
  button.active {
    background: var(--color-accent-muted);
  }

  .label {
    font-weight: 500;
  }

  .path {
    color: var(--color-text-muted);
    font-size: var(--text-xs);
    font-family: var(--font-mono);
  }

  .footer {
    display: flex;
    gap: var(--space-4);
    padding: var(--space-2) var(--space-4);
    border-top: 1px solid var(--color-border);
    font-size: var(--text-xs);
    color: var(--color-text-muted);
  }

  .footer kbd {
    font-family: var(--font-mono);
    margin-right: var(--space-1);
  }
</style>
