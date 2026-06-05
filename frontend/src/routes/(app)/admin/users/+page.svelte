<script lang="ts">
  import { onMount } from "svelte";
  import AdminGate from "$lib/admin/AdminGate.svelte";
  import { hasAdminKey } from "$lib/admin/adminKey";
  import { listAdminUsers, type AdminUser } from "$lib/admin/api";

  let unlocked = $state(hasAdminKey());
  let users = $state<AdminUser[]>([]);
  let busy = $state(false);
  let error = $state("");

  async function load() {
    if (!hasAdminKey()) return;
    busy = true;
    error = "";
    try {
      users = await listAdminUsers();
    } catch (e) {
      error = e instanceof Error ? e.message : "Falha ao carregar";
    } finally {
      busy = false;
    }
  }

  onMount(load);

  function onGateChange() {
    unlocked = hasAdminKey();
    if (unlocked) load();
  }
</script>

<svelte:head>
  <title>Utilizadores — Admin</title>
</svelte:head>

<section class="page">
  <a href="/admin" class="back">← Administração</a>
  <h1>Utilizadores</h1>

  <AdminGate onUnlocked={onGateChange} />

  {#if unlocked}
    {#if error}
      <p class="error" role="alert">{error}</p>
    {/if}
    {#if busy}
      <p class="muted">A carregar…</p>
    {:else if users.length === 0}
      <p class="muted">Nenhum utilizador registado.</p>
    {:else}
      <ul class="list">
        {#each users as u (u.ID)}
          <li>
            <a href="/admin/users/{u.ID}">
              <strong>{u.Email}</strong>
              <span class="meta">{u.ID.slice(0, 8)}… · {new Date(u.CreatedAt).toLocaleDateString("pt-PT")}</span>
            </a>
          </li>
        {/each}
      </ul>
    {/if}
  {/if}
</section>

<style>
  .page {
    max-width: 42rem;
  }

  .back {
    display: inline-block;
    margin-bottom: var(--space-4);
    color: var(--color-link);
    text-decoration: none;
    font-size: var(--text-sm);
  }

  h1 {
    margin: 0 0 var(--space-6);
    font-family: var(--font-display);
    font-size: var(--text-2xl);
  }

  .list {
    list-style: none;
    margin: 0;
    padding: 0;
    border: 1px solid var(--color-border);
    border-radius: var(--radius-md);
    overflow: hidden;
  }

  .list a {
    display: block;
    padding: var(--space-3) var(--space-4);
    text-decoration: none;
    color: inherit;
    border-bottom: 1px solid var(--color-border);
  }

  .list li:last-child a {
    border-bottom: none;
  }

  .list a:hover {
    background: var(--color-bg-surface);
  }

  .meta {
    display: block;
    font-size: var(--text-xs);
    color: var(--color-text-muted);
    font-family: var(--font-mono);
    margin-top: var(--space-1);
  }

  .error {
    padding: var(--space-3);
    margin-bottom: var(--space-4);
    border-radius: var(--radius-sm);
    background: color-mix(in srgb, var(--color-danger) 12%, transparent);
    color: var(--color-danger);
    font-size: var(--text-sm);
  }

  .muted {
    color: var(--color-text-muted);
    font-size: var(--text-sm);
  }
</style>
