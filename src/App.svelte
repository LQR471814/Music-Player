<script lang="ts">
import { IndexClient } from "./proto/api.client";
import { GrpcWebFetchTransport } from "@protobuf-ts/grpcweb-transport"
import type { Album } from "./proto/data";
import { Empty } from "./proto/api";

let albums: Album[] = []

const transport = new GrpcWebFetchTransport({
  baseUrl: "http://localhost:8000",
});
const client = new IndexClient(transport);
client.index(Empty).responses.onMessage((m) => {
  for (const update of m.updates) {
    switch (update.payload.oneofKind) {
      case "album":
        const a = update.payload.album
        console.info(`album: ${a.albumArtist} - ${a.title} | ${Object.keys(a.tracks).length} tracks`)
    }
  }
})
</script>

<main class="flex h-full w-full items-center justify-center" />

<style>
:global(body) {
  margin: 0;
  height: 100vh;
}
</style>
