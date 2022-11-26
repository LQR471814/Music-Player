<script lang="ts">
import { classList, sortAlphabetically } from "@web-std/common/src/general";
import { staggeredFly } from "@web-std/svelte-common/src/general";
import { store } from "~/store/state";
import AlbumChip from "../parts/AlbumChip.svelte";

export let className = "";

const albums = store.select((s) => s.albums);
$: chips = sortAlphabetically(Object.values($albums), (a) => a.title);
$: stagger = staggeredFly(Object.values($albums).length, {});
</script>

<div class={classList("flex-1", className)}>
  <div class="grid-container p-8">
    {#each chips as a, i (i)}
      <AlbumChip album={a} flyParams={stagger(i)} />
    {/each}
  </div>
</div>

<style>
.grid-container {
  display: grid;
  gap: 1rem;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  height: fit-content;
}
</style>
