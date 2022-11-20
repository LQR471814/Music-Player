<script lang="ts">
import { getContext } from "svelte";
import { Context, iconKey } from "./icon-context";

const context = getContext<Context>(iconKey);

export let className = context.className;
export let width = context.width;
export let height = context.height;
export let x = context.x;
export let y = context.y;
</script>

<svg
  class={className}
  viewBox="0 0 24 24"
  style={`width: ${width}px; height: ${height}px;`}
  {x}
  {y}
>
  <path fill="none" d="M0 0h24v24H0z"></path><path d="M12 10.586l4.95-4.95 1.414 1.414-4.95 4.95 4.95 4.95-1.414 1.414-4.95-4.95-4.95 4.95-1.414-1.414 4.95-4.95-4.95-4.95L7.05 5.636z"></path>
</svg>
