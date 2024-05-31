export function centToNormal(x) {
    return (x/100).toFixed(2)
}

export function normalToCent(x) {
    return Math.round(x*100)
}