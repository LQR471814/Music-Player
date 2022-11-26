import ColorThief, { RGBColor } from "colorthief"
import { fit, height, position, processImage, width } from "@web-std/common/src/images";

import { db } from "./db";
import { imageStore } from "@web-std/store/src/image"
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
    borderOpacity: 0.2,
}

export const light: Theme = {
    wallpaper: {
        type: "color",
        content: "white"
    },
    colors: {
        primary: "black",
        secondary: "black",
        background: "gray",
    },
    effects
}

export const dark: Theme = {
    wallpaper: {
        type: "color",
        content: "#363636"
    },
    colors: {
        primary: "white",
        secondary: "white",
        background: "#616161",
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

    const themeColor = document.querySelector("meta[name='theme-color']")
    if (themeColor) {
        (themeColor as HTMLMetaElement).content = theme.colors.background
    }
}

export const inferTheme = async (wallpaperId: number): Promise<Theme> => {
    const wallpaper = await db.wallpapers.get(wallpaperId)
    if (!wallpaper) {
        return light
    }

    const fullsizeURL = imageStore.fetch(wallpaper.fullsize)
    const referenceURL = await processImage(
        wallpaper.paletteReference,
        [
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
            },
        ],
    )
    if (!referenceURL) {
        return light
    }

    const img = new Image()
    await new Promise(r => {
        img.onload = r
        img.src = imageStore.fetch(referenceURL[0])
    })

    const colorThief = new ColorThief()
    const colors = colorThief.getPalette(img, 3)

    imageStore.release(referenceURL[0])

    return {
        colors: {
            primary: colorToString(colors[1]),
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
