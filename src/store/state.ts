import { fluxStore } from "@web-std/flux"
import type { Album, Track } from "~/proto/data";
import { BatchedUpdate, Action } from "~/proto/api"
import { orderedTracks } from "~/common/utils";
import { clamp, withoutElement } from "@web-std/common/src/general";

export type State = {
    playlist: {
        tracks: (Track & { album: string })[],
        position: number
        playing: boolean
    },
    albums: { [key: string]: Album }
    selectedAlbum?: string
    playerBlurred: boolean
}

const initialState: State = {
    albums: {},
    playlist: {
        tracks: [],
        position: 0,
        playing: false
    },
    playerBlurred: false
}

export const store = fluxStore(initialState, {
    batchedUpdates: (s, batch: BatchedUpdate) => {
        for (const update of batch.updates) {
            switch (update.payload.oneofKind) {
                case "album":
                    switch (update.action) {
                        case Action.ADD:
                        case Action.OVERRIDE:
                            s.albums[update.payload.album.id] = update.payload.album
                            break
                        case Action.REMOVE:
                            delete s.albums[update.payload.album.id]
                            break
                    }
                    break
                case "track":
                    const albumId = update.payload.track.albumId
                    const track = update.payload.track.track
                    if (!track) {
                        continue
                    }
                    switch (update.action) {
                        case Action.ADD:
                        case Action.OVERRIDE:
                            s.albums[albumId].tracks[track.id] = track
                            break
                        case Action.REMOVE:
                            delete s.albums[albumId].tracks[track.id]
                            break
                    }
                    break
            }
        }
    },
    selectAlbum: (s, album: string | undefined) => {
        if (s.selectedAlbum === album) {
            s.selectedAlbum = undefined
            return
        }
        s.selectedAlbum = album
    },
    playTrack: (s, track: number) => {
        s.playlist.position = clamp(track, 0, s.playlist.tracks.length-1)
        s.playlist.playing = true
    },
    removeTrack: (s, track: number) => {
        s.playlist.position = 0
        s.playlist.playing = false
        s.playlist.tracks = withoutElement(s.playlist.tracks, track)
    },
    togglePlaying: (s, playing?: boolean) => {
        if (playing) {
            s.playlist.playing = playing
            return
        }
        s.playlist.playing = !s.playlist.playing
    },
    overridePlaylist: (s, album: string) => {
        s.playlist.tracks = orderedTracks(s.albums[album]).map(o => ({ ...o, album }))
        s.playlist.playing = true
    },
    addToPlaylist: (s, album: string, track?: string) => {
        if (track) {
            s.playlist.tracks.push({ ...s.albums[album].tracks[track], album })
            return
        }
        for (const t of orderedTracks(s.albums[album])) {
            s.playlist.tracks.push({ ...t, album })
        }
    },
    playerBlurred: (s, value: boolean) => {
        s.playerBlurred = value
    }
})

export const currentAlbumSelector = (s: State) => s.selectedAlbum ?
    s.albums[s.selectedAlbum] : undefined
