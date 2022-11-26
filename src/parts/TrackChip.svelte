<script lang="ts">
import { twMerge } from "tailwind-merge";
import Label from "@web-std/form/src/Label.svelte";
import { fly, FlyParams } from "svelte/transition";
import type { Track } from "~/proto/data";

export let className = "";
export let track: Track;
export let index: number;
export let noDisc = false;
export let selected = false;
export let flyParams: FlyParams;
</script>

<button
  class={twMerge(
    "flex gap-4 p-3 rounded-2xl items-center interactive",
    selected ? "border-2 border-primary" : "",
    className
  )}
  transition:fly={flyParams}
  on:click
>
  <div class="flex flex-col items-start justify-center flex-1">
    <Label
      preset="h3"
      className={twMerge(
        "transition-all",
        selected ? "font-semibold" : "font-normal"
      )}
    >
      {!noDisc ? `${index + 1}. ${track.title}` : track.title}
    </Label>
    {#if track.artist.trim() !== ""}
      <Label preset="h4" className="transition-all font-normal">
        {track.artist}
      </Label>
    {/if}
  </div>
</button>
