import type { Album, MediaItem, Movie } from "./types";

export function title(s: string): string {
    return s.split(/\s/).map(s => s.charAt(0).toUpperCase() + s.slice(1)).join(' ')
}

export function isAlbum(item: MediaItem): item is MediaItem & { data: Album } {
    return item.type === "album"
}

export function isMovie(item: MediaItem): item is MediaItem & { data: Movie } {
    return item.type === "movie"
}
