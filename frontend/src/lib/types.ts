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
    type: string
    external_id: string
    data: Movie | Album
}

export interface ListItem {
    external_id: string
    type: string
}

export interface List {
    title: string
    created_by: string
    items: MediaItem[]
}

export interface ListPayload {
    title: string
    created_by: string
    items: ListItem[]
}
