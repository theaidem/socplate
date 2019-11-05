import { Room } from './room'

export interface User {
    id: number
    name: string
    avatar: string
    room?: Room["id"]
}
