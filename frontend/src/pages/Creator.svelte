<script>
    import MyButton from "@/lib/MyButton.svelte";
    import MyCard from "@/lib/MyCard.svelte";
    import MyInput from "@/lib/MyInput.svelte";
    import MySelector from "@/lib/MySelector.svelte";
    import { NFCWriteURL } from "@/lib/nfc";
    import { createAccount } from "@/lib/v1";
    import QRCode from "qrcode";

    let createJumper = false,
        showQRCode = false,
        accountRole = "normal";

    let result,
        nfcResult = "",
        qrCodeReady = false;

    function handleClear() {
        result = null;
        qrCodeReady = false;
        nfcResult = "";
    }

    function handleCreate() {
        createAccount(accountRole, createJumper)
            .then((r) => {
                result = r;
                alert(`账号创建成功`);

                if (result["link"] && showQRCode) {
                    QRCode.toCanvas(
                        document.getElementById("qrcode"),
                        result["link"],
                        (err) => {
                            if (err) {
                                console.log(err);
                            } else {
                                qrCodeReady = true;
                            }
                        },
                    );
                }
            })
            .catch((r) => {
                alert(r);
            });
    }

    async function handleWrite() {
        if (result && result["link"]) {
            try {
                nfcResult += `正在写入 Tag: URL=${result["link"]}\n`;
                await NFCWriteURL(result["link"]);
                nfcResult += `写入完成，无事发生\n`;
            } catch (err) {
                nfcResult += `错误: ${err}\n`;
            }
        } else {
            nfcResult += `没有需要写入 Tag 的 URL\n`;
        }
    }
</script>

<MyCard>
    <h1 class="font-bold text-xl">账户创建器</h1>
    <MyInput label="创建Jumper" type="checkbox" bind:value={createJumper}
    ></MyInput>
    <MyInput label="显示二维码" type="checkbox" bind:value={showQRCode}
    ></MyInput>
    <MySelector label="账户类型" bind:value={accountRole}>
        <option value="normal">普通</option>
        <option value="merchant">商户</option>
        <option value="operator">OP</option>
    </MySelector>
    <div>
        <MyButton primary on:click={handleCreate}>Just Do It!</MyButton>
        <MyButton on:click={handleWrite}>Write NFC</MyButton>
        <MyButton on:click={handleClear}>CLR</MyButton>
    </div>
    {#if nfcResult.length > 0}
        <h2>NFC运行日志</h2>
        <pre class="text-wrap break-all">{nfcResult}</pre>
    {/if}

    <canvas id="qrcode" hidden={!qrCodeReady}></canvas>

    {#if result != null}
        <h2>执行结果</h2>
        <pre class="text-wrap break-all">{JSON.stringify(result, null, 2)}</pre>
    {/if}
</MyCard>
