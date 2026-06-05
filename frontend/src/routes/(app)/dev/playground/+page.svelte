<script lang="ts">
  import PlaygroundDev from "$lib/playground/PlaygroundDev.svelte";
  import MotionDemo from "$lib/motion/MotionDemo.svelte";
  import { PageShell } from "$lib/ui";
</script>

<svelte:head>
  <title>Cofre playground — AegisPass</title>
</svelte:head>

<PageShell
  title="Cofre playground"
  taskId="VAULT-019"
  description="CRUD end-to-end legacy para testar o fluxo do cofre em desenvolvimento."
>
  <PlaygroundDev />
  <MotionDemo />
</PageShell>
