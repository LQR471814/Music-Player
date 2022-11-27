<script lang="ts">
import Label from "@web-std/form/src/Label.svelte";
import Menu from "@web-std/wrappers/src/Menu.svelte";
import AlbumIcon from "~/parts/AlbumIcon.svelte";
import PencilLine from "~/icons/PencilLine.svelte";
import DeleteBinLine from "~/icons/DeleteBinLine.svelte";
import PillLabel from "~/parts/PillLabel.svelte";

import { fly } from "svelte/transition";
import { classList } from "@web-std/common/src/general";
import { albumKey, receive, send } from "../common";
import { currentAlbumSelector, store } from "~/store/state";
import { layoutType, LayoutTypes } from "~/store/layout";

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
</script>

{#if $selectedAlbum}
  <!-- svelte-ignore a11y-click-events-have-key-events -->
  <div
    class={classList(
      "flex items-center gap-4 sm:gap-8 backdrop-blur-md p-6 sm:p-8 rounded-3xl",
      "border border-primary-clear border-opacity-20"
    )}
    transition:fly={{ y: 20 }}
    on:click
  >
    <div
      class="w-fit h-fit"
      in:receive={{ key: coverKey }}
      out:send={{ key: coverKey }}
    >
      <Menu
        side={$layoutType === LayoutTypes.PHONE ? "right" : "bottom"}
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
          size={$layoutType === LayoutTypes.PHONE ? 75 : 150}
          className={classList(
            "shadow-2xl",
            $layoutType === LayoutTypes.PHONE ? "rounded-2xl" : "rounded-3xl"
          )}
        />
      </Menu>
    </div>
    <div
      class="flex flex-col gap-1 sm:gap-2 justify-center flex-1 w-full"
      in:receive={{ key: titleKey }}
      out:send={{ key: titleKey }}
    >
      <Label
        className="font-semibold"
        preset={$layoutType === LayoutTypes.PHONE ? "h2" : "h1"}
        noMargin={$layoutType === LayoutTypes.PHONE}
      >
        {$selectedAlbum.title}
      </Label>
      {#if $selectedAlbum.albumArtist.trim() !== ""}
        <Label
          className="font-normal"
          preset={$layoutType === LayoutTypes.PHONE ? "h3" : "h2"}
        >
          {$selectedAlbum.albumArtist}
        </Label>
      {/if}
      {#if Object.keys($selectedAlbum.tracks).length === 1}
        <div
          class="w-fit h-fit mt-2"
          in:receive={{ key: singleKey }}
          out:send={{ key: singleKey }}
        >
          <PillLabel scale={$layoutType === LayoutTypes.PHONE ? "small" : "large"}>single</PillLabel>
        </div>
      {/if}
    </div>
  </div>
{/if}
