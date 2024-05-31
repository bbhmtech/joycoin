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
        getQuickAction()
            .then((r) => {
                quickActionSelected = r["action"];
                if (quickActionSelected == "quickPay") {
                    quickPayAmount = centToNormal(r["int64_value_1"]);
                    quickPayMessage = r["string_value_1"];
                } else {
                    (quickPayAmount = null), (quickPayMessage = null);
                }
            })
            .catch((r) => {
                console.log(r);
                quickActionSelected = "null";
            });
    });
    function handleSave() {
        postAccount(0, nickname, passcode);
        if (quickActionSelected == "quickPay") {
            setQuickPay(normalToCent(quickPayAmount), quickPayMessage, false);
        } else {
            clearQuickAction();
        }
        console.log(quickPayMessage);
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

    {#if role == "merchant" || role == "operator"}
        <MySelector
            label="[Staff] 本机快捷操作"
            bind:value={quickActionSelected}
        >
            <option value="null">无</option>
            <option value="quickPay">快捷支付</option>
        </MySelector>
    {/if}

    {#if quickActionSelected == "quickPay"}
        <MyInput
            type="number"
            label="快捷支付 - 金额"
            hint="Hint: 正数表示向对方支付"
            bind:value={quickPayAmount}
            placeholder="0.00"
        >
            <div class="pointer-events-none flex items-center">
                <span>🎲</span>
            </div>
        </MyInput>
        <MyInput
            type="text"
            label="快捷支付 - 描述"
            bind:value={quickPayMessage}
            placeholder="说点什么..."
        ></MyInput>
    {/if}

    <div class="flex justify-around">
        <MyButton on:click={handleSave} primary={true}>保存</MyButton>
        <MyButton on:click={handleBack}>返回</MyButton>
    </div>
</MyCard>
