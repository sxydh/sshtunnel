export function ifEqual(a: any, b: any, exclusive?: string[]): boolean {
    if (a === b) {
        return true
    }
    const keysA = Object.keys(a)
    const keysB = Object.keys(b)
    if (keysA.length !== keysB.length) {
        return false
    }
    for (let key of keysA) {
        if (exclusive && exclusive.includes(key)) {
            continue
        }
        if (!keysB.includes(key) || a[key] != b[key]) {
            return false
        }
    }
    return true
}