// const baseURL = "http://localhost:8080/_/v1"
const baseURL = "/_/v1"

async function get(path, params) {
    let searchParams = new URLSearchParams(params)
    let url = baseURL + path
    if (searchParams.size > 0) {
        url += '?' + searchParams.toString()
    }
    return await fetch(baseURL + path, {
        credentials: 'include'
    }).then(async r => {
        if (r.ok) {
            return await r.json()
        } else {
            return Promise.reject(await r.text())
        }
    });
}

async function post(path, params) {
    return await fetch(baseURL + path, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json;charset=utf-8'
        },
        credentials: 'include',
        body: JSON.stringify(params)
    }).then(async r => {
        if (r.ok) {
            return await r.json()
        } else {
            return Promise.reject(await r.text())
        }
    });
}

export async function getAccount(id) {
    return await get('/account/' + id)
}

export async function postAccount(id, nickname, passcode) {
    return await post('/account/' + id, {
        nickname, passcode
    })
}

export async function createAccount(role, createJumper) {
    return await post('/account/create', {
        'role': role,
        "create_jumper": createJumper
    })
}

export async function activateAccount(id, nickname, passcode) {
    return await post(`/account/${id}/activate`, {
        "nickname": nickname,
        "passcode": passcode
    })
}

export async function setQuickPay(centAmount, message, isTemporary) {
    return await post(`/quickaction`, {
        "action": "quickPay",
        "cent_amount": Number(centAmount),
        "message": String(message),
        "temporary": Boolean(isTemporary)
    })
}

export async function clearQuickAction() {
    return await post(`/quickaction`, {
        "action": "null"
    })
}

export async function getQuickAction() {
    return await get("/quickaction")
}

export async function listJumpers() {
    return await get(`/jumper`)
}

export async function listTransactions() {
    return await get('/transaction')
}