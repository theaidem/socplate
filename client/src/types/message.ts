type OutcomingEvent =
    | "create_room"
    | "join_room"
    | "leave_room"
    | "chat_message"

type IncomingEvent =
    | "user"
    | "online"
    | "room_list"
    | "start_room"
    | "update_room"
    | "chat_messages"

type Event = IncomingEvent | OutcomingEvent

interface Message<EventType extends Event> {
    event: EventType
    payload: any
}

export type IncomingMessage = Message<IncomingEvent>

export type OutcomingMessage = Message<OutcomingEvent>
