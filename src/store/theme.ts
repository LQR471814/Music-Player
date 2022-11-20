import { imageStore } from "@web-std/store/src/image"
import { db } from "./db";
import ColorThief, { RGBColor } from "colorthief"
import { fit, height, position, processImage, width } from "@web-std/common/src/images";
import { parseToRgba } from "color2k"

type Theme = {
    wallpaper: {
        type: "image" | "color"
        content: string
    }
    colors: {
        primary: string
        secondary: string
        background: string
    }
    effects: {
        backgroundOpacity: number
        borderOpacity: number
    }
}

const colorToString = (color: RGBColor) => {
    return `rgb(${color.toString()})`
}

export const effects: Theme["effects"] = {
    backgroundOpacity: 0.1,
    borderOpacity: 0.3,
}

export const light: Theme = {
    wallpaper: {
        type: "color",
        content: "white"
    },
    colors: {
        primary: "black",
        secondary: "black",
        background: "#e5e7eb",
    },
    effects
}

export const dark: Theme = {
    wallpaper: {
        type: "color",
        content: "gray"
    },
    colors: {
        primary: "white",
        secondary: "#7ba7c4",
        background: "#e5e7eb",
    },
    effects
}

export const loadTheme = (theme: Theme) => {
    const root = document.querySelector(":root") as HTMLElement | null
    if (!root) {
        return
    }

    for (const key in theme.colors) {
        const colorValue = parseToRgba(
            theme.colors[key as keyof Theme["colors"]]
        ).slice(0, 3).toString()

        root.style.setProperty(
            `--${key}`, colorValue
        )
    }

    root.style.setProperty(
        "--border-opacity",
        theme.effects.borderOpacity.toString(),
    )
    root.style.setProperty(
        "--bg-opacity",
        theme.effects.backgroundOpacity.toString(),
    )

    root.style.setProperty(
        "--wallpaper",
        theme.wallpaper.type === "image" ?
            `url(${theme.wallpaper.content})` :
            theme.wallpaper.content
    )
}

export const inferTheme = async (wallpaperId: number): Promise<Theme> => {
    const wallpaper = await db.wallpapers.get(wallpaperId)
    if (!wallpaper) {
        return light
    }

    const fullsizeURL = imageStore.fetch(wallpaper.fullsize)
    const referenceURL = await processImage(
        wallpaper.paletteReference,
        ({ context, image, canvas }) => {
            const fitted = position([0.5, 0.5], fit("contain", {
                min: [0, 0],
                max: [window.innerWidth, window.innerHeight],
            }, {
                min: [0, 0],
                max: [image.width, image.height],
            }), {
                min: [0, 0],
                max: [image.width, image.height]
            })
            const w = width(fitted)
            const h = height(fitted)

            canvas.width = w
            canvas.height = h
            context.drawImage(
                image, fitted.min[0], fitted.min[1], w, h,
                0, 0, w, h
            )
        }
    )
    if (!referenceURL) {
        return light
    }

    const img = new Image()
    await new Promise(r => {
        img.onload = r
        img.src = referenceURL
    })

    const colorThief = new ColorThief()
    const colors = colorThief.getPalette(img, 3)

    console.log(colors)

    return {
        colors: {
            primary: colorToString(colors[2]),
            secondary: colorToString(colors[2]),
            background: colorToString(colors[0]),
        },
        wallpaper: {
            type: "image",
            content: fullsizeURL,
        },
        effects,
    }
}
