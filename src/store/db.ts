import { Dexie, Table } from "dexie"
import { ImageSource, processImage, resize } from "@web-std/common/src/images"

const defaultsMap = ["1.jpg", "2.jpg"]

export type Wallpaper = {
    id?: number
    fullsize: Uint8Array
    paletteReference: Uint8Array
}

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
        const resizeImage = (size: number, blur?: string) => (
            { image, context, canvas }: {
                image: HTMLImageElement,
                context: CanvasRenderingContext2D,
                canvas: HTMLCanvasElement
            }
        ) => {
            const aspectRatio = image.height / image.width
            const width = (2 * size) / (aspectRatio + 1)

            canvas.width = width
            canvas.height = width * aspectRatio

            context.filter = `blur(${blur})`
            context.drawImage(image, 0, 0, width, width * aspectRatio)
        }

        const buffers = await processImage(source, [
            resizeImage(256),
            resizeImage(1280, "4px"),
        ])
        if (!buffers) {
            return null
        }

        return {
            fullsize: buffers[1],
            paletteReference: buffers[0],
        }
    }
}

export const db = new DB()
