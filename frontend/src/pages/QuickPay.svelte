<script>
    import MyButton from "@/lib/MyButton.svelte";
    import MyCard from "@/lib/MyCard.svelte";
    import MyInput from "@/lib/MyInput.svelte";
    import { centToNormal, normalToCent } from "@/lib/conv";
    import { clearQuickAction, getQuickAction, setQuickPay } from "@/lib/v1";
    import { onMount } from "svelte";

    let quickPaySelection = 0,
        quickPayAmount,
        quickPayMessage,
        quickPayRepeatable;

    onMount(() => {
        getQuickAction()
            .then((r) => {
                quickPayAmount = centToNormal(r["int64_value_1"]);
                quickPayMessage = r["string_value_1"];
                quickPayRepeatable = !r["temporary"]
                if (r["action"] == "quickPay") {
                    quickPaySelection = Number(quickPayAmount) <= 0 ? 1 : 2;
                } else {
                    quickPaySelection = 0
                }
            })
            .catch((r) => {
                quickPaySelection = 0;
            });
    });

    async function handleSave() {
        try {
            switch (quickPaySelection) {
                case 0:
                    await clearQuickAction();
                    break;
                case 1:
                    await setQuickPay(
                        normalToCent(-quickPayAmount),
                        quickPayMessage,
                        !quickPayRepeatable,
                    );
                    break;
                case 2:
                    await setQuickPay(
                        normalToCent(quickPayAmount),
                        quickPayMessage,
                        !quickPayRepeatable,
                    );
                    break;
                default:
                    break;
            }
            alert("已保存，1分钟贴近标签内有效");
        } catch (error) {
            console.log(error);
            alert(`错误: ${error}`);
        }
    }
    function handleBack() {
        history.back();
    }
</script>

<MyCard>
    <h1 class="font-bold text-xl">收付款</h1>
    <MyInput
        type="radio"
        label="选择"
        options={[
            { label: "取消", value: 0 },
            { label: "收", value: 1 },
            { label: "付", value: 2 },
        ]}
        bind:value={quickPaySelection}
    ></MyInput>
    <MyInput
        type="number"
        label="金额"
        bind:value={quickPayAmount}
        placeholder="0.00"
    >
        <div class="pointer-events-none flex items-center">
            <span>🎲</span>
        </div>
    </MyInput>
    
    <MyInput
        type="checkbox"
        label="是否多次收款"
        bind:value={quickPayRepeatable}
    ></MyInput>

    <MyInput
        type="text"
        label="附言"
        bind:value={quickPayMessage}
        placeholder="说点什么..."
    ></MyInput>
    <div class="flex justify-around">
        <MyButton on:click={handleSave} primary={true}>保存</MyButton>
        <MyButton on:click={handleBack}>返回</MyButton>
    </div>
</MyCard>
