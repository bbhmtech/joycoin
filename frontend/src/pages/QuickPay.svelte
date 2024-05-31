<script>
    import MyButton from "@/lib/MyButton.svelte";
    import MyCard from "@/lib/MyCard.svelte";
    import MyInput from "@/lib/MyInput.svelte";
    import { normalToCent } from "@/lib/conv";
    import { clearQuickAction, setQuickPay } from "@/lib/v1";

    let quickPaySelection = 0,
        quickPayAmount,
        quickPayMessage;
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
                        false,
                    );
                    break;
                case 2:
                    await setQuickPay(
                        normalToCent(quickPayAmount),
                        quickPayMessage,
                        false,
                    );
                    break;
                default:
                    break;
            }
            alert("已保存，1分钟贴近标签内有效");
        } catch (error) {
            console.log(error)
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
