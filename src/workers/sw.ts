/// <reference no-default-lib="true"/>
/// <reference lib="esnext" />
/// <reference lib="webworker" />
//@ts-ignore
const sw = self as unknown as ServiceWorkerGlobalScope & typeof globalThis;

sw.addEventListener("install", (e) => {
    console.log("service worker installed")
})
