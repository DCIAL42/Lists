export type TrackingStatus = "backlog" | "paused" | "done" | "none"
export type MediaType = "album" | "movie" | "game" | "show"

export interface Movie {
    title: string
    cover: string
}

export interface Album {
    title: string
    artist: string
    cover: string
}

export interface MediaItem {
    type: MediaType
    external_id: string
    data: Movie | Album
    tracking: {
        id?: number
        status?: TrackingStatus
    }
}

export interface ListItem {
    external_id: string
    type: MediaType
}

export interface List {
    id: number
    title: string
    created_by: string
    items?: MediaItem[]
}

export interface ListPayload {
    title: string
    created_by: string
    items: ListItem[]
}

export interface TrackingPayload {
    external_id: string
    status: TrackingStatus
    type: MediaType
}
