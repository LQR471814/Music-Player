<script lang="ts">
import { GrpcWebFetchTransport } from "@protobuf-ts/grpcweb-transport";
import { APIClient } from "~/proto/api.client";
import { Empty } from "~/proto/api";
import { apiLocation } from "./common/api";

import { store } from "./store/state";
import { setContext } from "svelte";
import { iconKey, Context as IconContext } from "./icons/icon-context";

import Library from "~/components/Library.svelte";
import { inferTheme, light, loadTheme } from "./store/theme";
import Playlist from "./components/Playlist.svelte";
import Player from "./components/Player.svelte";
import AlbumOverlay from "./components/AlbumOverlay.svelte";

import DynamicHeader, {
  Col as Column,
  Head as Header,
} from "@web-std/wrappers/src/DynamicHeader.svelte";
import { db } from "./store/db";

setContext<IconContext>(iconKey, {
  className: "w-14 h-14 svg-shadow",
});

const transport = new GrpcWebFetchTransport({
  baseUrl: apiLocation,
});

loadTheme(light)
db.loadDefaultWallpapers().then(() => {
  inferTheme(1).then((t) => {
    loadTheme(t);
  });
});

const client = new APIClient(transport);
client.index(Empty).responses.onMessage((m) => {
  store.actions.batchedUpdates(m);
});

const registration = navigator.serviceWorker.register("./workers/sw", {
  scope: "./workers/sw",
});

const playlistEmpty = store.select((s) => s.playlist.tracks.length === 0);
</script>

<main class="w-full h-full">
  <DynamicHeader>
    {#if !$playlistEmpty}
      <Column className="">
        <Playlist />
      </Column>
    {/if}
    <Column>
      <Library />
    </Column>
    <Header
      className="z-20 transition-all border-b border-transparent"
      coverClass="backdrop-blur-sm border-primary"
    >
      <Player />
    </Header>
  </DynamicHeader>
  <AlbumOverlay />
</main>
