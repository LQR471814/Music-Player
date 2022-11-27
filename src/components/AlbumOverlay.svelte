<script lang="ts">
import { currentAlbumSelector, store } from "~/store/state";

import IconButton from "~/parts/IconButton.svelte";

import TrackChip from "~/parts/TrackChip.svelte";
import PlayLine from "~/icons/PlayLine.svelte";
import PlayListAddLine from "~/icons/PlayListAddLine.svelte";
import { classList } from "@web-std/common/src/general";
import { useClose } from "@web-std/svelte-common/src/hooks";
import AlbumCard from "./AlbumOverlay/AlbumCard.svelte";
import { orderedTracks } from "~/common/utils";
import { staggeredFly } from "@web-std/svelte-common/src/general";
import { layoutType, LayoutTypes } from "~/store/layout";
import { clickOutside } from "@web-std/svelte-common/src/actions";

const selectedAlbum = store.select(currentAlbumSelector);

let albumId: string;
$: {
  if ($selectedAlbum) {
    albumId = $selectedAlbum.id;
  }
}

useClose(() => {
  store.actions.selectAlbum(undefined);
});
</script>

{#if $selectedAlbum}
  {@const tracks = orderedTracks($selectedAlbum)}
  {@const stagger = staggeredFly(tracks.length, { x: 0, y: 50 })}

  <!-- <div
    class="absolute w-[80%] h-[80%] rounded-3xl bg-white blur-xl"
    transition:fade={{ duration: 250 }}
  /> -->
  <!-- <div
    class="absolute w-full h-full rounded-3xl backdrop-blur-md"
    transition:fade={{ duration: 250 }}
  /> -->

  {#if $layoutType === LayoutTypes.DESKTOP}
    <!-- svelte-ignore a11y-click-events-have-key-events -->
    <div
      class={classList(
        "flex flex-row gap-8 justify-center",
        "overflow-y-auto absolute overlay z-30"
      )}
      use:clickOutside={{
        axis: (e) => e.children[0],
        callback: () => {
          store.actions.selectAlbum(undefined);
        },
      }}
    >
      <!-- on:click={() => {
        store.actions.selectAlbum(undefined);
      }} -->
      <!-- Album Info -->
      <div class="flex items-center h-full sticky top-0 z-40 max-w-xl">
        <AlbumCard />
      </div>
      <!-- Tracks -->
      <div class="flex h-full">
        <div class="flex flex-col gap-4 my-auto ml-0 py-16 pb-32">
          {#each tracks as t, i}
            <TrackChip
              className="max-w-md"
              track={t}
              index={i}
              flyParams={stagger(i)}
              on:click={(e) => {
                e.stopPropagation();
                store.actions.addToPlaylist(albumId, t.id);
              }}
            />
          {/each}
        </div>
      </div>
    </div>
  {:else}
    <!-- svelte-ignore a11y-click-events-have-key-events -->
    <div
      class={classList(
        "flex flex-col gap-8 items-center p-8",
        "overflow-y-auto absolute overlay z-30"
      )}
      on:click={() => {
        store.actions.selectAlbum(undefined);
      }}
    >
      <!-- Album Info -->
      <div
        class="flex items-center h-min sticky top-0 z-40 max-w-xl mt-32"
        on:click={(e) => e.stopPropagation()}
      >
        <AlbumCard />
      </div>
      <!-- Tracks -->
      <div class="flex flex-col gap-4 pb-36">
        {#each tracks as t, i}
          <TrackChip
            className="max-w-md"
            track={t}
            index={i}
            flyParams={stagger(i)}
            on:click={(e) => {
              e.stopPropagation();
              store.actions.addToPlaylist(albumId, t.id);
            }}
          />
        {/each}
      </div>
    </div>
  {/if}
  <!-- Controls -->
  <div class="fixed bottom-10 center-x flex gap-4 w-fit z-40">
    <IconButton
      icon={PlayLine}
      size={32}
      className="p-4"
      flyParams={{ y: 20 }}
      on:click={(e) => {
        store.actions.overridePlaylist(albumId);
        store.actions.selectAlbum(undefined);
      }}
    />
    <IconButton
      icon={PlayListAddLine}
      size={32}
      className="p-4"
      flyParams={{ y: 10 }}
      on:click={(e) => {
        store.actions.addToPlaylist(albumId);
      }}
    />
  </div>
{/if}
