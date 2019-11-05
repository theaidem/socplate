import { User } from "./user"
import { Room } from "./room"
import { Chat } from "./chat"
import { OutcomingMessage } from "./message"

type Stage = "lobby" | "room" | "settings" | "error"

export type DimmerStage = "load" | "error" | null

interface Dimming {
    active: boolean
    stage: DimmerStage | null
}

export interface CtxState {
    isInitialized: boolean
    dimming: Dimming
    user?: User | null
    stage: Stage
    room?: Room | null
    rooms: Room[]
    chat: Chat
    online: number
}

export interface CtxActions {
    sendMessage: (msg: OutcomingMessage) => void
}

interface Ctx extends CtxState, CtxActions { }

export interface CtxProps {
    ctx: Ctx
}
