export function title(s: string): string {
    return s.split(/\s/).map(s => s.charAt(0).toUpperCase() + s.slice(1)).join(' ')
}
