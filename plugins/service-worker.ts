import { rollup } from 'rollup'
import { nodeResolve } from '@rollup/plugin-node-resolve'
// @ts-ignore
import rollupPluginTypescript from 'rollup-plugin-typescript'

export default () => ({
  name: 'service-worker',
  writeBundle: async () => {
    const bundle = await rollup({
      input: 'src/workers/sw.ts',
      plugins: [rollupPluginTypescript(), nodeResolve()],
    })
    await bundle.write({
      file: 'src/workers/sw.js',
      format: 'es',
    })
    await bundle.close()
  }
})
