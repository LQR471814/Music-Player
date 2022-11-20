<script lang="ts">
import { fly } from "svelte/transition";
import { Label } from "@web-std/form";

import type { Album } from "~/proto/data";
import { classList } from "@web-std/common/src/general";
import { store } from "~/store/state";

import AlbumIcon from "./AlbumIcon.svelte";
import { albumKey, receive, send } from "~/components/common";
import PillLabel from "./PillLabel.svelte";

const selectedAlbum = store.select((s) => s.selectedAlbum);

export let album: Album;
export let index: number = 0;
</script>

<button
  class={classList(
    "flex gap-4 p-3 rounded-2xl items-center interactive",
    $selectedAlbum === album.id ? "border-secondary shadow-xl" : ""
  )}
  in:fly={{ delay: 35 * (index + 1), y: -10 }}
  on:click={() => {
    store.actions.selectAlbum(album.id);
  }}
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
      <Label preset="h3" className="font-semibold">{album.title}</Label>
      {#if album.albumArtist.trim() !== ""}
        <Label preset="h4" className="font-normal" noMargin>{album.albumArtist}</Label>
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
  {/if}
</button>
