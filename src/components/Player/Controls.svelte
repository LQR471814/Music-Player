<script lang="ts">
import Range from "@web-std/form/src/Range.svelte";
import PauseLine from "~/icons/PauseLine.svelte";
import PlayLine from "~/icons/PlayLine.svelte";
import SkipBackLine from "~/icons/SkipBackLine.svelte";
import SkipForwardLine from "~/icons/SkipForwardLine.svelte";
import ArrowLeftLine from "~/icons/ArrowLeftLine.svelte";
import IconButton from "~/parts/IconButton.svelte";

import { apiLocation } from "~/common/api";
import { currentPlayingAlbum, currentPlayingTrack, store } from "~/store/state";
import { controlPlay, safePadding } from "@web-std/svelte-common/src/actions";
import { imageStore } from "@web-std/store/src/image";
import { useKey } from "@web-std/svelte-common/src/hooks";
import Label from "@web-std/form/src/Label.svelte";
import { carryOver } from "@web-std/common/src/general";
import PlayListLine from "~/icons/PlayListLine.svelte";

export let showPlaylist = false;

const playlist = store.select((s) => s.playlist);
const playlistHidden = store.select((s) => s.playlist.hidden)
const playlistEmpty = store.select((s) => s.playlist.tracks.length === 0)
const album = store.select(currentPlayingAlbum);
const track = store.select(currentPlayingTrack);
const playing = store.select((s) => s.playlist.playing);

const changeTrack = (position: number) => {
  _elapsed = 0;
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
  changeTrack($playlist.position + 1);
});
navigator.mediaSession.setActionHandler("previoustrack", () => {
  changeTrack($playlist.position - 1);
});

useKey(" ", () => store.actions.togglePlaying());
useKey("P", () => changeTrack($playlist.position - 1));
useKey("N", () => changeTrack($playlist.position + 1));

let _dragging = false;
let _elapsed = 0;
let _player: HTMLAudioElement;

$: timestamp = (() => {
  if ($playlist.tracks.length === 0) {
    _elapsed = 0;
    return "-:--";
  }
  const carried = carryOver(_elapsed * (_player?.duration ?? 0), [
    60 * 60,
    60,
    1,
  ]);
  if (carried[0] > 0) {
    return carried.map((c) => c.toString().padStart(2, "0")).join(":");
  }
  return `${carried[1]}:${carried[2].toString().padStart(2, "0")}`;
})();

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
      if (!_player || _player.seeking || _dragging) return;
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
  class="flex flex-col gap-4 flex-1 px-8 py-6 phone:pt-4 sm:py-6"
  use:safePadding={{ sides: ["bottom"] }}
>
  <Range
    min={0}
    max={1}
    disabled={!$album}
    showLabel={false}
    outlineOpacity={0.3}
    trackClass="bg-transparent border border-primary-clear h-2"
    trackProgressClass="bg-primary h-3"
    thumbClass="hidden"
    bind:value={_elapsed}
    on:drag={(value) => {
      _elapsed = value.detail;
      _dragging = true;
    }}
    on:dragEnd={(value) => {
      if (!_player) return;
      _player.fastSeek(value.detail.value * _player.duration);
      _dragging = false;
    }}
  />
  <div class="flex gap-4 justify-center items-center relative">
    <Label className="absolute left-0" preset="h2">{timestamp}</Label>
    <IconButton
      icon={SkipBackLine}
      disabled={$playlist.position === 0}
      on:click={() => changeTrack($playlist.position - 1)}
      flyParams={{ duration: 0, delay: 0 }}
    />
    <IconButton
      icon={!$playing ? PlayLine : PauseLine}
      on:click={() => store.actions.togglePlaying()}
      disabled={$playlist.tracks.length === 0}
      flyParams={{ duration: 0, delay: 0 }}
    />
    <IconButton
      icon={SkipForwardLine}
      disabled={$playlist.position === $playlist.tracks.length - 1 ||
        $playlist.tracks.length === 0}
      on:click={() => changeTrack($playlist.position + 1)}
      flyParams={{ duration: 0, delay: 0 }}
    />
    {#if showPlaylist}
      <IconButton
        className="absolute right-0"
        icon={$playlistHidden || $playlistEmpty ? PlayListLine : ArrowLeftLine}
        disabled={$playlist.tracks.length === 0}
        flyParams={{ duration: 0, delay: 0 }}
        on:click={store.actions.togglePlaylist}
      />
    {/if}
  </div>
</div>
