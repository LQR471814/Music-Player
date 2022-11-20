<script lang="ts">
import { currentAlbumSelector, store } from "~/store/state";
import { albumKey, receive, send } from "./common";

import { Label } from "@web-std/form";
import { Menu } from "@web-std/wrappers";
import AlbumIcon from "~/parts/AlbumIcon.svelte";
import IconButton from "~/parts/IconButton.svelte";

import PencilLine from "~/icons/PencilLine.svelte";
import DeleteBinLine from "~/icons/DeleteBinLine.svelte";
import PillLabel from "~/parts/PillLabel.svelte";
import TrackChip from "~/parts/TrackChip.svelte";
import PlayLine from "~/icons/PlayLine.svelte";
import PlayListAddLine from "~/icons/PlayListAddLine.svelte";
import { fly } from "svelte/transition";
import { classList } from "@web-std/common/src/general";
import { clickOutside } from "@web-std/svelte-common/src/actions";
import { useClose } from "@web-std/svelte-common/src/hooks";

const selectedAlbum = store.select(currentAlbumSelector);

let albumId: string;
$: {
  if ($selectedAlbum) {
    albumId = $selectedAlbum.id;
  }
}

$: coverKey = albumKey(albumId, "cover");
$: titleKey = albumKey(albumId, "title");
$: singleKey = albumKey(albumId, "single");

useClose(() => {
  store.actions.selectAlbum(undefined);
});
</script>

{#if $selectedAlbum}
  <div class="absolute flex gap-8 overlay overflow-y-auto z-30">
    <!-- <div
      class="absolute w-[80%] h-[80%] rounded-3xl bg-white blur-xl"
      transition:fade={{ duration: 250 }}
    /> -->
    <!-- <div
      class="absolute w-full h-full rounded-3xl backdrop-blur-md"
      transition:fade={{ duration: 250 }}
    /> -->
    <!-- Album Info -->
    <div
      class="flex justify-end items-center flex-1 sticky top-0"
      use:clickOutside={{
        callback: () => {
          store.actions.selectAlbum(undefined);
        },
      }}
    >
      <div
        class={classList(
          "flex flex-col items-end gap-8 backdrop-blur-md p-8 rounded-3xl",
          "border border-primary-clear border-opacity-20"
        )}
        transition:fly={{ y: 20 }}
      >
        <div
          class="w-fit h-fit"
          in:receive={{ key: coverKey }}
          out:send={{ key: coverKey }}
        >
          <Menu
            side="left"
            options={[
              {
                title: "Upload",
                icon: PencilLine,
                onaction: () => {},
              },
              {
                title: "Clear",
                icon: DeleteBinLine,
                onaction: () => {},
              },
            ]}
          >
            <AlbumIcon
              album={$selectedAlbum}
              size={216}
              className="shadow-2xl rounded-3xl"
            />
          </Menu>
        </div>
        <div
          class="flex flex-col gap-2 items-end justify-center flex-1 max-w-[300px]"
          in:receive={{ key: titleKey }}
          out:send={{ key: titleKey }}
        >
          <Label className="text-end font-semibold" preset="h1">
            {$selectedAlbum.title}
          </Label>
          {#if $selectedAlbum.albumArtist.trim() !== ""}
            <Label className="text-end font-normal" preset="h2">
              {$selectedAlbum.albumArtist}
            </Label>
          {/if}
        </div>
        {#if Object.keys($selectedAlbum.tracks).length === 1}
          <div
            class="w-fit h-fit"
            in:receive={{ key: singleKey }}
            out:send={{ key: singleKey }}
          >
            <PillLabel scale="large">single</PillLabel>
          </div>
        {/if}
      </div>
    </div>
    <!-- Tracks -->
    <div class="py-32 h-fit flex-1">
      <div class="flex flex-col gap-4">
        {#each Object.values($selectedAlbum.tracks) as t, i}
          <TrackChip
            className="max-w-md"
            track={t}
            index={i}
            on:click={(e) => {
              e.stopPropagation();
              store.actions.addToPlaylist(albumId, t.id);
            }}
          />
        {/each}
      </div>
    </div>
    <!-- Controls -->
    <div class="fixed bottom-10 center-x flex gap-4 w-fit">
      <IconButton
        className="p-4"
        flyParams={{ y: 20 }}
        on:click={(e) => {
          store.actions.overridePlaylist(albumId);
        }}
      >
        <PlayLine width={32} height={32} />
      </IconButton>
      <IconButton
        className="p-4"
        flyParams={{ y: 10 }}
        on:click={(e) => {
          store.actions.addToPlaylist(albumId);
        }}
      >
        <PlayListAddLine width={32} height={32} />
      </IconButton>
    </div>
  </div>
{/if}
