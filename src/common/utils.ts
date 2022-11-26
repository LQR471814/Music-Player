import type { Album } from "~/proto/data";
import type { Writable } from "svelte/store";

export function orderedTracks(a: Album) {
    return Object.values(a.tracks).sort((a, b) => a.disc - b.disc)
}

export type StoreValue<S> = S extends Writable<infer T> ? T : never
