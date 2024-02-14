export type Status = "completed" | "failed" | "running" | "queued"

export type Message = {
    id: string;
    completedAt?: string;
    restarts: number;
    createdAt: string;
    status: Status;
    mType: String
}

