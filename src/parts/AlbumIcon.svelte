<script lang="ts">
import type { Album } from "~/proto/data";
import { Image } from "@web-std/utility";
import AlbumLine from "~/icons/AlbumLine.svelte";
import AlbumFill from "~/icons/AlbumFill.svelte";

export let album: Album;
export let size = 56;
export let className = "shadow-lg rounded-xl"
</script>

{#if album.cover}
    <Image
        className={className}
        image={album.cover.data}
        alt={album.cover.description}
        minLength={size}
        maxLength={size}
    />
{:else if Object.keys(album.tracks).length === 1}
    <AlbumLine width={size} height={size} />
{:else}
    <AlbumFill width={size} height={size} />
{/if}
