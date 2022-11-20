<script lang="ts">
import AlbumFill from "~/icons/AlbumFill.svelte";
import PlayLine from "~/icons/PlayLine.svelte";
import SkipForwardLine from "~/icons/SkipForwardLine.svelte";
import SkipBackLine from "~/icons/SkipBackLine.svelte";
import AlbumIcon from "~/parts/AlbumIcon.svelte";
import IconButton from "~/parts/IconButton.svelte";
import { Label, Range } from "@web-std/form";
import type { Track } from "~/proto/data";
import { store, State } from "~/store/state";
import PauseLine from "~/icons/PauseLine.svelte";
import { apiLocation } from "~/common/api";

import { controlPlay } from "@web-std/svelte-common/src/actions";
import { useKey } from "@web-std/svelte-common/src/hooks";
import { classList } from "@web-std/common/src/general";
import { imageStore } from "@web-std/store/src/image";

const selectTrack: (s: State) => (Track & { album: string }) | undefined = (
  s
) => s.playlist.tracks[s.playlist.position];

const playlist = store.select((s) => {
  return s.playlist;
});
const album = store.select((s) => {
  const albumId = selectTrack(s)?.album;
  if (albumId) {
    return s.albums[albumId];
  }
});
const track = store.select(selectTrack);
const playing = store.select((s) => s.playlist.playing);
const blurred = store.select((s) => s.playerBlurred);

const changeTrack = (position: number) => {
  if (position >= $playlist.tracks.length) {
    store.actions.playTrack(0);
    return;
  }
  store.actions.playTrack(position);
};

navigator.mediaSession.setActionHandler("play", () => {
  store.actions.togglePlaying(true);
});
navigator.mediaSession.setActionHandler("pause", () => {
  store.actions.togglePlaying(false);
});
navigator.mediaSession.setActionHandler("nexttrack", () => {
  changeTrack($playlist.position + 1)
});
navigator.mediaSession.setActionHandler("previoustrack", () => {
  changeTrack($playlist.position - 1)
});

useKey(" ", () => store.actions.togglePlaying());
useKey("P", () => changeTrack($playlist.position - 1));
useKey("N", () => changeTrack($playlist.position + 1));


let _elapsed = 0;
let _player: HTMLAudioElement;

$: {
  if (_player) {
    if ($playing) {
      navigator.mediaSession.metadata = new MediaMetadata({
        title: $track?.title ?? "no title",
        album: $album?.title ?? "unknown album",
        artist: $track?.artist ?? "unknown artist",
        artwork: $album?.cover
          ? [
              {
                src: imageStore.fetch($album.cover.data),
                sizes: "96x96",
                type: "image/jpeg",
              },
            ]
          : undefined,
      });
    }
    _player.ontimeupdate = () => {
      if (!_player || _player.seeking) return;
      _elapsed = _player.currentTime / _player.duration;
    };
    _player.onended = () => {
      if (!_player) return;
      changeTrack($playlist.position + 1);
    };
  }
}
</script>

{#if $track}
  <audio
    bind:this={_player}
    autoplay
    src={(() => {
      const url = new URL(apiLocation);
      url.pathname = $track.path;
      return url.toString();
    })()}
    use:controlPlay={{ store: playing }}
  />
{/if}

<div
  class={classList(
    "flex items-center gap-6 p-8",
    $blurred ? "backdrop-blur-sm" : ""
  )}
>
  <div>
    {#if $album}
      <AlbumIcon size={80} album={$album} />
    {:else}
      <AlbumFill width={80} height={80} />
    {/if}
  </div>
  <div>
    <Label preset="h2" className="mb-1 max-w-sm line-clamp-2">
      {$track?.title ?? "not playing"}
    </Label>
    <Label preset="h3" className="max-w-sm line-clamp-1">
      {$track?.artist ?? "unknown artist"}
    </Label>
  </div>
  <div class="flex flex-col gap-4 flex-1 pl-4">
    <Range
      min={0}
      max={1}
      disabled={!$album}
      showLabel={false}
      outlineOpacity={0.3}
      trackClass="bg-transparent backdrop-blur-md border border-primary-clear"
      trackProgressClass="bg-primary h-[6px]"
      thumbClass="bg-primary border-none outline-primary"
      bind:value={_elapsed}
      on:dragEnd={(value) => {
        if (!_player) return;
        _player.fastSeek(value.detail.value * _player.duration);
      }}
    />
    <div class="flex gap-4 justify-center">
      <IconButton
        disabled={$playlist.position === 0}
        on:click={() => changeTrack($playlist.position - 1)}
      >
        <SkipBackLine width={24} height={24} />
      </IconButton>
      <IconButton
        on:click={() => store.actions.togglePlaying()}
        disabled={$playlist.tracks.length === 0}
      >
        {#if !$playing}
          <PlayLine width={24} height={24} />
        {:else}
          <PauseLine width={24} height={24} />
        {/if}
      </IconButton>
      <IconButton
        disabled={$playlist.position === $playlist.tracks.length - 1 ||
          $playlist.tracks.length === 0}
        on:click={() => changeTrack($playlist.position + 1)}
      >
        <SkipForwardLine width={24} height={24} />
      </IconButton>
    </div>
  </div>
</div>
