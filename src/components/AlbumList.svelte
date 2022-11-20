<script lang="ts">
//** UNUSED */

import { classList } from "@web-std/common/src/general";
import { rotate } from "@web-std/svelte-common/src/transitions";
import { store } from "~/store/state";
import { fly } from "svelte/transition";

import { Label } from "@web-std/form";
import AlbumIcon from "~/parts/AlbumIcon.svelte";
import TrackChip from "~/parts/TrackChip.svelte";
import Edit2Line from "~/icons/Edit2Line.svelte";

export let className = "";

const selectedAlbum = store.select((s) =>
  s.selectedAlbum ? s.albums[s.selectedAlbum] : undefined
);
</script>

{#if $selectedAlbum}
  <div
    class={classList("flex flex-col h-full w-1/3", className)}
    transition:fly={{ x: 10, duration: 200 }}
  >
    <div class="flex flex-col items-center gap-8 p-8">
      <AlbumIcon
        album={$selectedAlbum}
        size={168}
        className="shadow-2xl rounded-3xl"
      />
      <Label preset="h1" className="text-start">{$selectedAlbum.title}</Label>
    </div>
    <div class="flex flex-col gap-4 p-8 overflow-y-auto">
      {#each Object.values($selectedAlbum.tracks) as t, i}
        <TrackChip track={t} index={i} />
      {/each}
      <button
        class={classList(
          "absolute bottom-8 right-8 rounded-full p-3",
          "shadow-lg bg-white interactive"
        )}
        in:fly={{ delay: 200, x: 10 }}
      >
        <div in:rotate={{}}>
          <Edit2Line width={32} height={32} />
        </div>
      </button>
    </div>
  </div>
{/if}
