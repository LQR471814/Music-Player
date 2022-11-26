import { windowSize } from "@web-std/store/src/window";
import { writable } from "svelte/store";

export enum LayoutTypes {
    PHONE = 0,
    TABLET = 1,
    DESKTOP = 2,
}

export const layoutType = writable<LayoutTypes>(LayoutTypes.PHONE)
windowSize.subscribe((dimensions) => {
    layoutType.set(
        dimensions[0] <= 500 ?
            LayoutTypes.PHONE :
            dimensions[0] <= 1080 ?
                LayoutTypes.TABLET :
                LayoutTypes.DESKTOP
    )
});