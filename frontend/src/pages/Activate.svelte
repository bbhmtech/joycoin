<script>
    import MyButton from "@/lib/MyButton.svelte";
    import MyCard from "@/lib/MyCard.svelte";
    import MyInput from "@/lib/MyInput.svelte";
    import { activateAccount, getAccount } from "@/lib/v1";
    import { onMount } from "svelte";
    import { push, querystring } from "svelte-spa-router";
    let nickname = "",
        passcode = "";

    let qs = new URLSearchParams($querystring);
    let id = qs.get("id");
    let initial = Boolean(qs.get("initial"));
    onMount(() => {
        if (!qs.has("id")) {
            alert("[错误] ID 缺失");
        }
    });

    function handleClick() {
        activateAccount(Number(id), nickname, passcode).then((r) =>
            push("/account"), (r) => alert(r)
        );
    }
</script>

<MyCard>
        <h1>Welcome</h1>
        {#if initial}
            <p>你的账号需要激活</p>
            <MyInput
                type="text"
                label="昵称"
                bind:value={nickname}
                placeholder="说点什么..."
            ></MyInput>
        {:else}
            <p>请登录你的账号</p>
        {/if}
        <MyInput
            type="password"
            label="口令"
            bind:value={passcode}
            placeholder="1234"
        ></MyInput>
        <MyButton primary on:click={handleClick}
            >{initial ? "激活" : "登录"}</MyButton
        >
</MyCard>
