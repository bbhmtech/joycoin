<script>
    export let primary = false;
</script>

{#if primary}
    <button
        class="pointer-events-auto rounded-md bg-indigo-600 px-3 py-2 font-semibold leading-5 text-white hover:bg-indigo-500"
        on:click
    >
        <slot />
    </button>
{:else}
    <button
        class="pointer-events-auto rounded-md bg-slate-400 px-3 py-2 font-semibold leading-5 text-white hover:bg-slate-500"
        on:click
    >
    <slot />
    </button>
{/if}
