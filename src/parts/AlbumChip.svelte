<script lang="ts">
import { fly, FlyParams } from "svelte/transition";
import Label from "@web-std/form/src/Label.svelte";

import type { Album } from "~/proto/data";
import { classList, styleList } from "@web-std/common/src/general";
import { store } from "~/store/state";

import AlbumIcon from "./AlbumIcon.svelte";
import { albumKey, receive, send } from "~/components/common";
import PillLabel from "./PillLabel.svelte";
import { useResizeObserver } from "@web-std/svelte-common/src/hooks";

const selectedAlbum = store.select((s) => s.selectedAlbum);

export let album: Album;
export let flyParams: FlyParams;

let unselectedHeight = 0;
const attachObserver = useResizeObserver((e) => {
  if ($selectedAlbum !== album.id) {
    unselectedHeight = e[0].contentRect.height;
  }
});
</script>

<button
  class={classList(
    "w-full flex gap-4 p-3 rounded-2xl items-center interactive",
    $selectedAlbum === album.id ? "border-secondary shadow-xl" : ""
  )}
  on:click={() => {
    store.actions.selectAlbum(album.id);
  }}
  transition:fly={flyParams}
  use:attachObserver
>
  {#if album.id !== $selectedAlbum}
    <div
      in:receive={{ key: albumKey(album.id, "cover") }}
      out:send={{ key: albumKey(album.id, "cover") }}
    >
      <AlbumIcon {album} />
    </div>
    <div
      class="flex flex-col items-start justify-center flex-1"
      in:receive={{ key: albumKey(album.id, "title") }}
      out:send={{ key: albumKey(album.id, "title") }}
    >
      <Label preset="h3" className="font-normal">{album.title}</Label>
      {#if album.albumArtist.trim() !== ""}
        <Label preset="h4" className="font-normal" noMargin>
          {album.albumArtist}
        </Label>
      {/if}
    </div>
    {#if Object.keys(album.tracks).length === 1}
      <div
        class="w-fit h-fit"
        in:receive={{ key: albumKey(album.id, "single") }}
        out:send={{ key: albumKey(album.id, "single") }}
      >
        <PillLabel>single</PillLabel>
      </div>
    {/if}
  {:else}
    <!-- * to retain the chip's original height -->
    <div style={styleList({ height: `${unselectedHeight}px` })} />
  {/if}
</button>
