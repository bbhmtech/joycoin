<script>
    export let label,
        hint = "",
        value = null;
</script>

<label>
    {label}
    {#if hint}
        <p class="text-gray-600">{hint}</p>
    {/if}

    <div class="relative rounded-md shadow-sm">
        <select
            bind:value
            class="block w-full rounded-md border-0 py-1.5 pl-3 pr-20 text-gray-900 ring-1 ring-inset ring-gray-300 placeholder:text-gray-400 focus:ring-2 focus:ring-inset focus:ring-indigo-600"
        >
            <slot></slot>
        </select>
    </div>
</label>
