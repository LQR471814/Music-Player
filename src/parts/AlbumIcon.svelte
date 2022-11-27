<script lang="ts">
import type { Album } from "~/proto/data";
import AlbumLine from "~/icons/AlbumLine.svelte";
import AlbumFill from "~/icons/AlbumFill.svelte";
import { styleList } from "@web-std/common/src/general";
import { iconLocation } from "~/common/api";

export let album: Album;
export let size = 56;
export let className = "shadow-lg rounded-xl";

$: src = iconLocation(album);
</script>

{#if src && album.cover}
  <img
    class={className}
    {src}
    alt={album.cover.description}
    style={styleList({
      objectFit: "contain",
      width: `${size}px`,
    })}
  />
{:else if Object.keys(album.tracks).length === 1}
  <AlbumLine width={size} height={size} />
{:else}
  <AlbumFill width={size} height={size} />
{/if}
