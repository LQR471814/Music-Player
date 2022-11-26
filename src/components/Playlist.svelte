<script lang="ts">
import { classList } from "@web-std/common/src/general";
import { store } from "~/store/state";
import TrackChip from "~/parts/TrackChip.svelte";
import Menu from "@web-std/wrappers/src/Menu.svelte";
import PlayLine from "~/icons/PlayLine.svelte";
import CloseLine from "~/icons/CloseLine.svelte";
import { staggeredFly } from "@web-std/svelte-common/src/general";

const playlist = store.select((s) => s.playlist);
let selected: number | undefined;

$: stagger = staggeredFly($playlist.tracks.length, {});
</script>

<div
  class="phone:w-full phone:max-w-none md:max-w-xs md:min-w-[280px] lg:max-w-sm p-8 border-primary"
>
  <div class="flex flex-col gap-4">
    {#each $playlist.tracks as t, i}
      <Menu
        side="bottom"
        containerClass={classList("w-full", selected === i ? "z-10" : "")}
        selectedScaling="110%"
        menuOffset="10%"
        options={[
          {
            title: "Play",
            icon: PlayLine,
            onaction: () => {
              store.actions.playTrack(i);
            },
          },
          {
            title: "Remove",
            icon: CloseLine,
            onaction: () => {
              store.actions.removeTrack(i);
            },
          },
        ]}
        on:select={(e) => {
          const [isSelected, _] = e.detail;
          if (isSelected) {
            selected = i;
            return;
          }
          selected = undefined;
        }}
      >
        <TrackChip
          className="w-full"
          track={t}
          index={i}
          selected={i === $playlist.position}
          noDisc
          flyParams={stagger(i)}
        />
      </Menu>
    {/each}
  </div>
</div>
