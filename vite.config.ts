import { resolve } from "path"
import { defineConfig } from 'vite';
import { svelte } from '@sveltejs/vite-plugin-svelte';
import sveltePreprocess from 'svelte-preprocess';
import fs from 'fs';
import VitePWA from "vite-pwa"

import serviceWorker from "./plugins/service-worker"

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
      serviceWorker(),
      VitePWA(),
    ]
  }
})