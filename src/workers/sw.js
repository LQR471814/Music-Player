/// <reference no-default-lib="true"/>
/// <reference lib="esnext" />
/// <reference lib="webworker" />
//@ts-ignore
const sw = self;
sw.addEventListener("install", (e) => {
    console.log("service worker installed");
});
