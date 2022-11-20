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
  <path fill="none" d="M0 0h24v24H0z"></path><path d="M8 11.333l10.223-6.815a.5.5 0 0 1 .777.416v14.132a.5.5 0 0 1-.777.416L8 12.667V19a1 1 0 0 1-2 0V5a1 1 0 1 1 2 0v6.333zm9 4.93V7.737L10.606 12 17 16.263z"></path>
</svg>
