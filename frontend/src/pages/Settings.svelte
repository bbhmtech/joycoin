<script>
    import MyButton from "@/lib/MyButton.svelte";
    import MyCard from "@/lib/MyCard.svelte";
    import MyInput from "@/lib/MyInput.svelte";
    import MySelector from "@/lib/MySelector.svelte";
    import { centToNormal, normalToCent } from "@/lib/conv";
    import {
        clearQuickAction,
        getAccount,
        getQuickAction,
        postAccount,
        setQuickPay,
    } from "@/lib/v1";
    import { onMount } from "svelte";

    let quickActionSelected = "null",
        passcode = "",
        nickname = "",
        role = "";
    let quickPayMessage, quickPayAmount;

    onMount(() => {
        getAccount(0)
            .then((r) => {
                nickname = r["nickname"];
                role = r["role"];
            })
            .catch((r) => {
                alert(r);
            });
    });
    function handleSave() {
        postAccount(0, nickname, passcode);
        alert("如存");
    }
    function handleBack() {
        history.back();
    }
</script>

<MyCard>
    <h1 class="font-bold text-xl">账户设置</h1>

    <MyInput type="text" label="昵称" bind:value={nickname}></MyInput>
    <MyInput
        type="password"
        label="登录口令"
        bind:value={passcode}
        placeholder="留空就不修改了"
    ></MyInput>

    <div class="flex justify-around">
        <MyButton on:click={handleSave} primary={true}>保存</MyButton>
        <MyButton on:click={handleBack}>返回</MyButton>
    </div>
</MyCard>
