<script lang="ts">
import { GrpcWebFetchTransport } from "@protobuf-ts/grpcweb-transport";
import { APIClient } from "~/proto/api.client";
import { Empty } from "~/proto/api";
import { apiLocation } from "./common/api";

import { store } from "./store/state";
import { setContext } from "svelte";
import { iconKey, Context as IconContext } from "./icons/icon-context";
import { inferTheme, dark, light, loadTheme } from "./store/theme";
import { db } from "./store/db";
import { classList } from "@web-std/common/src/general";

import Library from "~/components/Library.svelte";
import Playlist from "./components/Playlist.svelte";
import PlayerControls from "./components/Player/Controls.svelte";
import PlayerInfo from "./components/Player/Info.svelte";
import AlbumOverlay from "./components/AlbumOverlay.svelte";

import DynamicHeader, {
  Col as Column,
  Head as Header,
  Foot as Footer,
} from "@web-std/wrappers/src/DynamicHeader.svelte";
import { windowSize } from "@web-std/store/src/window";
import { layoutType, LayoutTypes } from "./store/layout";
import type { StoreValue } from "./common/utils";

setContext<IconContext>(iconKey, {
  className: "w-14 h-14 svg-shadow",
});

const transport = new GrpcWebFetchTransport({
  baseUrl: apiLocation,
});

loadTheme(light);
// db.loadDefaultWallpapers().then(() => {
//   inferTheme(1).then((t) => {
//     loadTheme(t);
//   });
// });

const client = new APIClient(transport);
client.index(Empty).responses.onMessage((m) => {
  store.actions.batchedUpdates(m);
});

// const registration = navigator.serviceWorker.register("./workers/sw", {
//   scope: "./workers/sw",
// });

const playlistEmpty = store.select((s) => s.playlist.tracks.length === 0);
const playlistHidden = store.select((s) => s.playlist.hidden);

let layout: StoreValue<typeof layoutType>;
layoutType.subscribe((l) => (layout = l));
</script>

<main class="w-full h-full">
  <DynamicHeader className="gap-0">
    <!-- * Playlist and Library -->
    {#if layout >= LayoutTypes.TABLET}
      {#if !$playlistEmpty}
        {#if layout > LayoutTypes.PHONE}
          <Column className="">
            <Playlist />
          </Column>
        {/if}
        <Column className="" paddingClass="h-full">
          <div class="w-0 h-full border-l border-primary-clear" />
        </Column>
      {/if}
      <Column>
        <Library />
      </Column>
    {:else}
      <Column>
        <Library />
      </Column>
      {#if !$playlistEmpty && !$playlistHidden}
        <Column className="absolute top-0 backdrop-blur-sm transition-opacity w-full">
          <Playlist />
        </Column>
      {/if}
    {/if}

    <!-- * Player Info & Controls -->
    <Header
      className={classList(
        "z-20 transition-all border-b",
        !$playlistEmpty ? "border-primary-clear" : "border-transparent",
        layout <= LayoutTypes.TABLET ? "rounded-b-[2rem]" : ""
      )}
      coverClass="backdrop-blur-sm border-primary-clear"
    >
      {#if layout === LayoutTypes.DESKTOP}
        <div class="flex gap-8">
          <PlayerInfo />
          <PlayerControls />
        </div>
      {:else}
        <PlayerInfo />
      {/if}
    </Header>
    {#if layout <= LayoutTypes.TABLET}
      <Footer
        className={classList(
          "transition-all border-t",
          !$playlistEmpty ? "border-primary-clear" : "border-transparent",
          layout <= LayoutTypes.TABLET ? "rounded-t-[2rem]" : ""
        )}
        coverClass="backdrop-blur-sm border-primary-clear"
      >
        <PlayerControls showPlaylist={layout === LayoutTypes.PHONE} />
      </Footer>
    {/if}
  </DynamicHeader>
  <AlbumOverlay />
</main>
