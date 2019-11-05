import React, { MouseEvent } from "react"
import { Segment, List, Header, Button, Icon } from "semantic-ui-react"
import { CtxProps } from "../../types/context"
import { Room } from "../../types/room"
import Item from "./item"

export default class RoomList extends React.Component<CtxProps> {
    public render() {
        const { ctx } = this.props
        if (!ctx.user) return null

        return (
            <Segment className="col display-flex flex-one flex-direction-col">
                <Header as="h3" dividing={true} size="small">
                    <Icon name="users" />
                    Rooms
                </Header>
                <div className="room-list display-flex flex-one">
                    {ctx.rooms.length > 0 && (
                        <List
                            divided={true}
                            relaxed={true}
                            className="room-list-items display-flex"
                        >
                            {ctx.rooms.map((room: Room) => (
                                <Item
                                    key={room.id}
                                    sendMessage={ctx.sendMessage}
                                    room={room}
                                    user={ctx.user}
                                />
                            ))}
                        </List>
                    )}
                </div>
                <div>
                    {ctx.user.room === null && (
                        <Button
                            content="Create room"
                            size="mini"
                            icon={{ color: "red", name: "add" }}
                            onClick={this.createRoomHandle}
                        />
                    )}
                </div>
            </Segment>
        )
    }

    private createRoomHandle = (e: MouseEvent) => {
        const { sendMessage } = this.props.ctx
        e.preventDefault()
        sendMessage({
            event: "create_room",
            payload: null,
        })
    }
}
