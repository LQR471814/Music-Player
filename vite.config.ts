import { resolve } from "path"
import { defineConfig } from 'vite';
import { svelte } from '@sveltejs/vite-plugin-svelte';
import sveltePreprocess from 'svelte-preprocess';
import fs from 'fs';
import VitePWA from "vite-pwa"

// import serviceWorker from "./plugins/service-worker"

export default defineConfig(({ mode }) => {
  return {
    root: "src",
    publicDir: "../public",
    server: {
      port: 8080,
      https: {
        key: fs.readFileSync('./server/host-key.pem'),
        cert: fs.readFileSync('./server/host-crt.pem'),
      }
    },
    build: {
      outDir: "../build",
      sourcemap: true,
    },
    resolve: {
      alias: {
        "~": resolve(__dirname, "src")
      },
    },
    plugins: [
      svelte({
        // @ts-ignore
        preprocess: sveltePreprocess({
          postcss: true
        }),
      }),
      // serviceWorker(),
      VitePWA({
        // @ts-ignore
        registerType: "autoUpdate",
        includeAssets: [
          "assets/**/*",
          "backgrounds/*",
          "favicon.png",
          "index.html",
        ],
        manifest: {
          name: "music player",
          short_name: "music",
          description: "a web-based music player",
          theme_color: "#ffffff",
          icons: [
            {
              src: "icons/android-chrome-192x192.png",
              sizes: "192x192",
              type: "image/png",
            },
            {
              src: "icons/android-chrome-512x512.png",
              sizes: "512x512",
              type: "image/png",
            },
            {
              src: 'icons/android-chrome-512x512.png',
              sizes: '512x512',
              type: 'image/png',
              purpose: 'any maskable'
            }
          ]
        }
      }),
    ]
  }
})