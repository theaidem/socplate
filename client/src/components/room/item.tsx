import React, { MouseEvent, Component } from "react"
import { List, Label, Button } from "semantic-ui-react"
import { Room as IRoom } from "../../types/room"
import { User } from "../../types/user"
import { CtxActions } from "../../types/context"

export default class Item extends Component<{
    room: IRoom
    user?: User | null
    sendMessage: CtxActions["sendMessage"]
}> {
    public render() {
        const { room, user } = this.props
        if (!user) return null

        return (
            <List.Item className="display-flex">
                <List.Content className="room-list-item display-flex">
                    <div className="user-list display-flex">
                        {room.users.map((u: User) => (
                            <Label key={u.id} as="a" image={true}>
                                <img alt={u.name} src={u.avatar} />
                                {u.name}
                            </Label>
                        ))}
                    </div>
                    <div className="room-actions display-flex">
                        {user.room === room.id && (
                            <Button
                                content="leave"
                                size="mini"
                                icon={{ name: "sign-out" }}
                                onClick={this.leaveRoomHandle}
                            />
                        )}
                        {user.room === null && (
                            <Button
                                content="join"
                                size="mini"
                                icon={{ name: "sign-in" }}
                                onClick={(e: MouseEvent) =>
                                    this.joinRoomHandle(e, room.id)
                                }
                            />
                        )}
                    </div>
                </List.Content>
            </List.Item>
        )
    }

    private joinRoomHandle = (e: MouseEvent, roomID: string) => {
        const { sendMessage } = this.props
        e.preventDefault()
        sendMessage({
            event: "join_room",
            payload: roomID,
        })
    }

    private leaveRoomHandle = (e: MouseEvent) => {
        const { sendMessage } = this.props
        e.preventDefault()
        sendMessage({
            event: "leave_room",
            payload: null,
        })
    }
}
