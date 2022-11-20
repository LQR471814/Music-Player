import type { Album } from "~/proto/data";

export function orderedTracks(a: Album) {
    return Object.values(a.tracks).sort((a, b) => a.disc - b.disc)
}