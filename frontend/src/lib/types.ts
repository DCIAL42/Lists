export type TrackingStatus = "backlog" | "paused" | "done" | "none"
export type MediaType = "album" | "movie" | "game" | "show"

export type SearchResponse = {
    [K in MediaType]?: { next: string; items: MediaItem[] };
};

export interface Movie {
    popularity: number
}

export interface Album {
    artist: string
}

export interface TrackingItem {
    id: number
    status: TrackingStatus
}

export interface MediaItem {
    id: number
    type: MediaType
    title: string
    cover: string
    data: Movie | Album
    tracking?: TrackingItem
    rating: Rating
}

export interface ListItem {
    media_id: number
    type: MediaType
}

export interface ListMeta {
    id: number
    title: string
    created_by: string
    cover: string
    likes: number
}

export interface List extends ListMeta {
    items: MediaItem[]
}

export interface ListPayload {
    title: string
    items: ListItem[]
}

export interface TrackingPayload {
    media_id: number
    status: TrackingStatus
    type: MediaType
}

export interface UserResponse {
    id: string
    username: string
}

export interface ListsPreviewData {
    lists: ListMeta[];
    next: string;
    page: number;
    count: number;
}

export interface TrackingListData {
    items: MediaItem[];
    count: number;
}

export type ProfileData =
    | {
        self: true
        trackingData: TrackingListData;
        listsData: ListsPreviewData;
        userData: UserResponse;
    } | {
        self: false
        listsData: ListsPreviewData;
        userData: UserResponse;
        followData: Follow
    }

export interface Follow {
    id?: number
    followed: boolean
}

export interface Like {
    id?: number
    liked: boolean
}

export interface Rating {
    id?: number
    rating: number
}
