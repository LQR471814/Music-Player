import { quintOut } from 'svelte/easing';
import { crossfade } from "svelte/transition";

export const albumKey = (
    albumId: string,
    component: "cover" | "title" | "artist" | "single"
) => {
    return albumId + "-" + component
}

export const [send, receive] = crossfade({
    duration: 300,
    fallback(node, params) {
        const style = getComputedStyle(node);
        const transform = style.transform === 'none' ? '' : style.transform;

        return {
            duration: 600,
            easing: quintOut,
            css: t => `
                transform: ${transform} scale(${t});
                opacity: ${t}
            `
        };
    }
});
