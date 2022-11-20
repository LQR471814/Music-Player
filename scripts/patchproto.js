import { readdir, readFile, writeFile } from "fs/promises"
import { join } from "path"

const args = process.argv.slice(2)
const protoDir = args[0] ?? "src/proto"

const importResolver = /^import .+ from ".\/(.+)";?/gm

async function addExtensions(path) {
    const f = await readFile(path, {
        encoding: "utf8"
    })
    await writeFile(path, f.replace(importResolver, function(match, src) {
        return match.replace(`./${src}`, `./${src}.js`)
    }))
}

for (const directory of await readdir(protoDir)) {
    addExtensions(join(protoDir, directory))
}
