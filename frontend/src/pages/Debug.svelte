<script>
    import { activateAccount, setQuickPay, getAccount, listJumpers } from "@/lib/v1";

    let id = 0,
        passcode = "",
        nickname = "",
        balance = 0;

    async function handleActivateAccount() {
        await activateAccount(id, nickname, passcode);
    }

    async function handleUpdateAccount() {
        let ret = await getAccount(id);
        nickname = ret["Nickname"];
        balance = ret["CachedCentBalance"];
    }

    let qa_action = "quickPay",
        qa_message = "",
        qa_temporary = false,
        qa_cent_amount = 0;
    async function handleSetQuickAction() {
        await setQuickPay(qa_cent_amount, qa_message, qa_temporary);
    }

    async function handleListJumpers() {
        console.log(await listJumpers())
    }
</script>

<label>
    Account ID
    <input bind:value={id} />
</label>
<label>
    Passcode
    <input bind:value={passcode} />
</label>
<label>
    Nickname
    <input bind:value={nickname} />
</label>
<p>
    Balance={balance}
</p>
<button on:click={handleUpdateAccount}>Account Detail</button>
<button on:click={handleActivateAccount}>Account Activate</button>
<label>
    QuickAction
    <input bind:value={qa_action} />
</label>
<label>
    IsTemporary
    <input type="checkbox" bind:value={qa_temporary} />
</label>
<label>
    Message
    <input bind:value={qa_message} />
</label>
<label>
    CentAmount
    <input bind:value={qa_cent_amount} />
</label>
<button on:click={handleSetQuickAction}>Set QuickPay</button>
<button on:click={handleListJumpers}>List Jumpers</button>
