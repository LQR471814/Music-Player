<script lang="ts">
import {
  ComponentConstructorOptions,
  createEventDispatcher,
  SvelteComponentTyped,
} from "svelte";
import { fly, FlyParams } from "svelte/transition";
import { twMerge } from "tailwind-merge";

interface RequiredProps {
  width?: number;
  height?: number;
  className?: string;
}

interface ComponentConstructor {
  new (o: ComponentConstructorOptions): SvelteComponentTyped<RequiredProps>;
}

const dispatcher = createEventDispatcher<{ click: MouseEvent }>();

export let icon: ComponentConstructor;
export let size: number = 24;

export let className = "";
export let flyParams: FlyParams = { y: 10 };
export let disabled = false;
</script>

<button
  class={twMerge(
    "interactive p-2 rounded-full",
    disabled
      ? "hover:scale-100 !bg-background-clear !border-background-clear"
      : "hover:scale-110 active:scale-90",
    className
  )}
  transition:fly={flyParams}
  {disabled}
  on:click={(e) => {
    dispatcher("click", e);
  }}
>
  <svelte:component
    this={icon}
    width={size}
    height={size}
    className={disabled ? "fill-background" : ""}
  />
</button>
