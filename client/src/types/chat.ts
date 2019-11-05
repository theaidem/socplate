import { User } from "./user"

export interface ChatMessage {
    from: User
    message: string
    created: number
}

export type Chat = ChatMessage[]
