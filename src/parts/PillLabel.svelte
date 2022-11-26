<script lang="ts">
import Label from "@web-std/form/src/Label.svelte";
import { classList } from "@web-std/common/src/general";

export let scale: "small" | "large" = "small";
</script>

<div
  class={classList(
    "rounded-full border border-primary",
    scale === "small" ? "px-2 py-1" : "px-4 py-2"
  )}
>
  <Label preset={scale === "small" ? "h5" : "h3"} noBold><slot /></Label>
</div>
