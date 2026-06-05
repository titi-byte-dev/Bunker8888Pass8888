<script lang="ts">
  import { NAV_ICON_PATHS, type NavIconName } from "./navIcons";

  interface Props {
    name: NavIconName;
    size?: number;
  }

  let { name, size = 20 }: Props = $props();
</script>

<svg
  class="nav-icon"
  width={size}
  height={size}
  viewBox="0 0 24 24"
  aria-hidden="true"
  focusable="false"
>
  <path d={NAV_ICON_PATHS[name]} fill="currentColor" />
</svg>

<style>
  .nav-icon {
    flex-shrink: 0;
    display: block;
  }
</style>
