<script>
    import { centToNormal } from "@/lib/conv";
    import { activateAccount, getAccount, listTransactions } from "@/lib/v1";
    import { onMount } from "svelte";
    import { push } from "svelte-spa-router";

    let nickname = "无名客", centBalance = 0, aId;
    let transactions = [];
    onMount(() => {
        getAccount(0)
            .then((r) => {
                nickname = r["nickname"];
                aId = r["id"];
                centBalance = r["cached_cent_balance"];
            })
            .catch(() => {
                push("/unauthorized");
            });

        listTransactions().then((r) => {
            transactions = r;
        });
    });
</script>

<main class="h-min-screen flex flex-col justify-center items-center">
    <div
        class="h-fit w-10/12 p-4 rounded-md shadow-md bg-slate-200 flex flex-col gap-4"
    >
        <h1 class="font-lg">6.1 快乐, {nickname}!</h1>
        <h1 class="font-bold text-xl">欢乐豆账户 ({aId})</h1>
        <h2 class="w-fit self-center font-mono text-4xl">
            <span class="align-middle text-lg">🎲</span>{centToNormal(centBalance)}
        </h2>
        <div class="flex justify-around">
            <button class="p-2 rounded-lg bg-sky-400">向个人支付</button>
            <button class="p-2 rounded-lg bg-sky-400">打赏</button>
            <a class="p-2 rounded-lg bg-gray-400" href="#/settings">设置</a>
            <button></button>
        </div>
        <h2 class="font-semibold text-lg">交易记录</h2>
        <table class="table-auto">
            <thead class="text-left border-t-2 border-slate-400">
                <tr>
                    <th>金额</th>
                    <th>内容</th>
                    <th>时间</th>
                    <th>流向</th>
                </tr>
            </thead>
            <tbody class="border-y-2 border-slate-400">
                {#each transactions as t, i}
                    <tr>
                        <td>🎲{centToNormal(t.cent_amount)}</td>
                        <td>{t.message}</td>
                        <td>1</td>
                        <td>{t.from_account_id} -> {t.to_account_id}</td>
                    </tr>
                {/each}
            </tbody>
        </table>
    </div>
</main>
