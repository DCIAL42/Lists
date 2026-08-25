export type TrackingStatus = "backlog" | "paused" | "done" | "none"
export type MediaType = "album" | "movie" | "game" | "show"

export type SearchResponse = {
    [K in MediaType]?: { next: string; items: MediaItem[] };
};

export interface Movie {
    title: string
}

export interface Album {
    title: string
    artist: string
}

export interface TrackingItem {
    id: number
    status: TrackingStatus
}

export interface MediaItem {
    type: MediaType
    external_id: string
    cover: string
    data: Movie | Album
    tracking?: TrackingItem
}

export interface ListItem {
    external_id: string
    type: MediaType
}

export interface ListMeta {
    id: number
    title: string
    created_by: string
    cover: string
}

export interface List extends ListMeta {
    items: MediaItem[]
}

export interface ListPayload {
    title: string
    items: ListItem[]
}

export interface TrackingPayload {
    external_id: string
    status: TrackingStatus
    type: MediaType
}
