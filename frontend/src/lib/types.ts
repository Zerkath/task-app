export type Status = "completed" | "failed" | "running" | "queued"

export type Page = {
    data: Message[];
    page: number;
    count: number;
}

export type Message = {
    id: string;
    completedAt?: string;
    restarts: number;
    createdAt: string;
    status: Status;
    mType: String
}

export type PageSettings = {
    pageSize: number;
    currentPage: number;
    pages: number;
}

