<script lang="ts">
import Range from "@web-std/form/src/Range.svelte";
import PauseLine from "~/icons/PauseLine.svelte";
import PlayLine from "~/icons/PlayLine.svelte";
import SkipBackLine from "~/icons/SkipBackLine.svelte";
import SkipForwardLine from "~/icons/SkipForwardLine.svelte";
import ArrowLeftLine from "~/icons/ArrowLeftLine.svelte";
import IconButton from "~/parts/IconButton.svelte";

import { apiLocation, iconLocation } from "~/common/api";
import { currentPlayingAlbum, currentPlayingTrack, store } from "~/store/state";
import { controlPlay, safePadding } from "@web-std/svelte-common/src/actions";
import { imageStore } from "@web-std/store/src/image";
import { useKey } from "@web-std/svelte-common/src/hooks";
import Label from "@web-std/form/src/Label.svelte";
import { carryOver } from "@web-std/common/src/general";
import PlayListLine from "~/icons/PlayListLine.svelte";
import {
  fit,
  height,
  position,
  processImageRaw,
  width,
} from "@web-std/common/src/images";

export let showPlaylist = false;

const playlist = store.select((s) => s.playlist);
const playlistHidden = store.select((s) => s.playlist.hidden);
const playlistEmpty = store.select((s) => s.playlist.tracks.length === 0);
const album = store.select(currentPlayingAlbum);
const track = store.select(currentPlayingTrack);
const playing = store.select((s) => s.playlist.playing);

$: hasNext = $playlist.position < $playlist.tracks.length - 1;
$: hasPrevious = $playlist.position > 0;

const changeTrack = (position: number) => {
  _elapsed = 0;
  if (position >= $playlist.tracks.length) {
    store.actions.playTrack(0);
    return;
  }
  store.actions.playTrack(position);
};

$: {
  if (!$playlistEmpty) {
    navigator.mediaSession.setActionHandler("play", () => {
      store.actions.togglePlaying(true);
    });
    navigator.mediaSession.setActionHandler("pause", () => {
      store.actions.togglePlaying(false);
    });
  } else {
    navigator.mediaSession.setActionHandler("play", null);
    navigator.mediaSession.setActionHandler("pause", null);
  }
}

navigator.mediaSession.setActionHandler("nexttrack", () => {
  changeTrack($playlist.position + 1);
});
navigator.mediaSession.setActionHandler("previoustrack", () => {
  changeTrack($playlist.position - 1);
});

track.subscribe((t) => {
  if (!t) return;

  // const getArtwork = async () => {
  //   if (!$album?.cover) {
  //     return undefined;
  //   }

  //   const resized = await processImageRaw($album?.cover.data, [
  //     ({ image, context, canvas }) => {
  //       const coverWidth = 256;
  //       const fitted = position(
  //         [0.5, 0.5],
  //         fit(
  //           "contain",
  //           {
  //             min: [0, 0],
  //             max: [coverWidth, coverWidth],
  //           },
  //           {
  //             min: [0, 0],
  //             max: [image.width, image.height],
  //           }
  //         ),
  //         {
  //           min: [0, 0],
  //           max: [image.width, image.height],
  //         }
  //       );

  //       const w = width(fitted);
  //       const h = height(fitted);

  //       canvas.width = coverWidth;
  //       canvas.height = coverWidth;

  //       context.drawImage(
  //         image,
  //         fitted.min[0],
  //         fitted.min[1],
  //         w,
  //         h,
  //         0,
  //         0,
  //         coverWidth,
  //         coverWidth
  //       );
  //     },
  //   ]);
  //   if (!resized) {
  //     return undefined;
  //   }

  //   console.log(resized[0])

  //   return [
  //     {
  //       src: resized[0],
  //       sizes: "256x256",
  //       type: "image/jpeg",
  //     },
  //   ];
  //   // return [
  //   //   {
  //   //     src: "https://localhost:8080/icons/android-chrome-192x192.jpeg",
  //   //     sizes: "192x192",
  //   //     type: "image/jpeg",
  //   //   },
  //   // ];
  // };
  // getArtwork().then((a) => {
  const iconURL = $album ? iconLocation($album) : undefined;
  console.log(iconURL);
  navigator.mediaSession.metadata = new MediaMetadata({
    title: $track?.title ?? "no title",
    album: $album?.title ?? "unknown album",
    artist: $track?.artist ?? "unknown artist",
    artwork: iconURL ? [{ src: iconURL }] : undefined,
  });
  // });
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
      url.pathname = `/audio/${$track.path}`;
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
      disabled={!hasPrevious || $playlistEmpty}
      on:click={() => changeTrack($playlist.position - 1)}
      flyParams={{ duration: 0, delay: 0 }}
    />
    <IconButton
      icon={!$playing ? PlayLine : PauseLine}
      on:click={() => store.actions.togglePlaying()}
      disabled={$playlistEmpty}
      flyParams={{ duration: 0, delay: 0 }}
    />
    <IconButton
      icon={SkipForwardLine}
      disabled={!hasNext || $playlistEmpty}
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
