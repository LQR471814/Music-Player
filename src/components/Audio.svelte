<script lang="ts">
export let src = "";
export let player: HTMLAudioElement;
export let loaded = false;

$: {
  loaded = false;
  src;
}
</script>

<audio
  bind:this={player}
  on:load={() => {
    loaded = true;
  }}
  {src}
/>
