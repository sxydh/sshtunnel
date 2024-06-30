export function stringify(p: any, exclusive?: string[]): string {
    let str = JSON.stringify(p)
    if (exclusive) {
        const json = JSON.parse(str)
        for (let e of exclusive) {
            delete json[e]
        }
        str = JSON.stringify(json)
    }
    return str
}