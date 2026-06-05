<script lang="ts">
  // Tokens partilhados com a app (gerado por scripts/sync-tokens.mjs).
  import "$lib/tokens.generated.css";
  import "../app.css";
  import type { Snippet } from "svelte";

  let { children }: { children: Snippet } = $props();
</script>

{@render children()}
