<script lang="ts">
  import AuthShell from "$lib/auth/AuthShell.svelte";
</script>

<svelte:head>
  <title>AegisPass — Bem-vindo</title>
</svelte:head>

<AuthShell
  title="Bem-vindo"
  subtitle="Gestão de identidade Zero-Knowledge para equipas BYOD."
>
  <div class="welcome">
    <a class="primary" href="/auth/login">Iniciar sessão</a>
    <a class="secondary" href="/auth/register">Criar conta</a>
    <a class="link" href="/auth/recovery">Recuperar acesso</a>
  </div>
</AuthShell>

<style>
  .welcome {
    display: flex;
    flex-direction: column;
    gap: var(--space-2);
  }

  .welcome a {
    display: block;
    padding: var(--space-3) var(--space-4);
    border-radius: var(--radius-sm);
    text-align: center;
    text-decoration: none;
    font-size: var(--text-sm);
    font-weight: 600;
  }

  .primary {
    background: var(--color-accent);
    color: var(--color-accent-fg);
  }

  .secondary {
    background: var(--color-accent-muted);
    color: var(--color-text);
  }

  .link {
    margin-top: var(--space-2);
    color: var(--color-link);
    font-weight: 500;
  }
</style>
