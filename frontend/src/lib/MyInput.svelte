<script>
    export let label,
        hint = "",
        type,
        value = null,
        placeholder = "";
</script>

<label>
    {label}
    {#if hint}
        <p class="text-gray-600">{hint}</p>
    {/if}
    <div
        class="flex gap-1 rounded-md shadow-sm border-0 py-1.5 px-3 bg-white text-gray-900 ring-1 ring-inset ring-gray-300 placeholder:text-gray-400 focus-within:ring-2 focus-within:ring-inset focus-within:ring-indigo-600"
    >
        <slot></slot>
        {#if type === "password"}
            <input
                type="password"
                class="flex-grow outline-none"
                bind:value
                {placeholder}
            />
        {:else if type === "email"}
            <input
                type="email"
                class="flex-grow outline-none"
                bind:value
                {placeholder}
            />
        {:else if type === "number"}
            <input
                type="number"
                class="flex-grow outline-none"
                bind:value
                {placeholder}
            />
        {:else if type === "date"}
            <input
                type="date"
                class="flex-grow outline-none"
                bind:value
                {placeholder}
            />
        {:else if type == "checkbox"}
            <input
                type="checkbox"
                class="flex-grow outline-none"
                bind:checked={value}
            />
        {:else}
            <input
                type="text"
                class="flex-grow outline-none"
                bind:value
                {placeholder}
            />
        {/if}
    </div>
</label>
