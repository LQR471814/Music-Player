import type { Album } from "~/proto/data"

// @ts-ignore
export const apiLocation = import.meta.env.PROD ?
    // window.location.origin :
    `https://${window.location.hostname}:6325` :
    `https://${window.location.hostname}:6325`

export const iconLocation = (album: Album) => {
    if (!album?.cover) {
        return undefined
    }
    return `${apiLocation}/${album.cover.url}`
}