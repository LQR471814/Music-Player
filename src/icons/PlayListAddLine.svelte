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
  <path fill="none" d="M0 0h24v24H0z"></path><path d="M2 18h10v2H2v-2zm0-7h20v2H2v-2zm0-7h20v2H2V4zm16 14v-3h2v3h3v2h-3v3h-2v-3h-3v-2h3z"></path>
</svg>
