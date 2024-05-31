<script>
    import MyCard from "@/lib/MyCard.svelte";
    import { centToNormal } from "@/lib/conv";
    import { getTransaction } from "@/lib/v1";
    import { onMount } from "svelte";
    import { querystring } from "svelte-spa-router";

    let qs = new URLSearchParams($querystring);
    let id = qs.get("id");
    let trans,
        success,
        reason,
        timeout = 15;
    onMount(async () => {
        if (!qs.has("id")) {
            alert("[错误] ID 缺失");
        } else {
            await getTransaction(Number(id))
                .then((r) => {
                    trans = r;
                    success = true;
                    reason = null;
                })
                .catch((r) => {
                    trans = null;
                    success = false;
                    reason = r;
                })
                .finally(() => {
                    setInterval(() => {
                        timeout--;
                        if (timeout <= 0) {
                            window.close();
                        }
                    }, 1000);
                });
        }
    });
</script>

<MyCard>
    {#if success}
        <h1 class="text-xl font-bold text-center text-green-400">交易完成！</h1>

        <h2 class="w-fit self-center font-mono text-4xl">
            <span class="align-middle text-lg">🎲</span>{centToNormal(
                trans["cent_amount"],
            )}
        </h2>
    {:else}
        <h1 class="text-xl font-bold text-center text-red-400">交易失败</h1>
        <h2>错误信息</h2>
        <pre>{reason}</pre>
    {/if}
    {#if timeout < 15}
        <h2>{timeout}s 后自动关闭</h2>
    {/if}
</MyCard>
