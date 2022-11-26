<script lang="ts">
import Label from "@web-std/form/src/Label.svelte";
import AlbumFill from "~/icons/AlbumFill.svelte";
import AlbumIcon from "~/parts/AlbumIcon.svelte";

import { currentPlayingAlbum, currentPlayingTrack, store } from "~/store/state";

const album = store.select(currentPlayingAlbum);
const track = store.select(currentPlayingTrack);
</script>

<div class="flex gap-8 items-center px-6 py-4 sm:py-6">
  <div>
    {#if $album}
      <AlbumIcon size={80} album={$album} />
    {:else}
      <AlbumFill width={80} height={80} />
    {/if}
  </div>
  <div>
    <Label preset="h2" className="mb-1 max-w-sm line-clamp-2 font-semibold">
      {$track?.title ?? "not playing"}
    </Label>
    <Label preset="h3" className="max-w-sm line-clamp-1" noBold>
      {$track?.artist ?? "unknown artist"}
    </Label>
  </div>
</div>
