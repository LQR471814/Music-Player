import { Dexie, Table } from "dexie"
import { ImageSource, processImage } from "@web-std/common/src/images"

const defaultsMap = ["1.jpg", "2.jpg"]

export type Wallpaper = {
    id?: number
    fullsize: Uint8Array
    paletteReference: string
}

const resizeSize = 256

class DB extends Dexie {
    wallpapers!: Table<Wallpaper>

    constructor() {
        super('wallpapers')
        this.version(1).stores({
            wallpapers: '++id'
        })
    }

    async loadDefaultWallpapers() {
        const wallpapers: Wallpaper[] = []
        let id = 0
        for (const src of defaultsMap) {
            if (await this.wallpapers.get({ id })) {
                continue
            }
            const wallpaper = await this.processImage(`backgrounds/${src}`)
            if (wallpaper) {
                wallpaper.id = id
                wallpapers.push(wallpaper)
            }
            id++
        }
        await this.wallpapers.bulkPut(wallpapers)
    }

    async processImage(source: ImageSource): Promise<Wallpaper | null> {
        let imageBuffer = new Uint8Array()

        const dataURL = await processImage(source, ({ buffer, image, context, canvas }) => {
            imageBuffer = buffer
            const aspectRatio = image.height / image.width
            const width = (2 * resizeSize) / (aspectRatio + 1)

            canvas.width = width
            canvas.height = width * aspectRatio

            context.drawImage(image, 0, 0, width, width * aspectRatio)
        })
        if (!dataURL) {
            return null
        }

        return {
            fullsize: imageBuffer,
            paletteReference: dataURL,
        }
    }
}

export const db = new DB()
