import React, { MouseEvent, Component } from "react"
import { Header, Segment, Icon, Label } from "semantic-ui-react"
import Chat from "../chat"
import { CtxProps } from "../../types/context"
import { User } from "../../types/user"

class Room extends Component<CtxProps> {
    public render() {
        const { ctx } = this.props
        if (!ctx.room) return null

        return (
            <div className="display-flex flex-one">
                <div className="col display-flex flex-one flex-direction-col">
                    <Segment className="col display-flex flex-one flex-direction-col">
                        <Header as="h3" dividing={true} size="small">
                            <Icon name="users" />
                            Users
                        </Header>
                        <div className="user-list display-flex flex-one">
                            {ctx.room.users.map((u: User) => (
                                <Label as="div" image={true} key={u.id}>
                                    <img alt={u.name} src={u.avatar} />
                                    {u.name}
                                    {u.id === (ctx.user ? ctx.user.id : 0) && (
                                        <Label.Detail
                                            as="a"
                                            onClick={this.leveRoomHandle}
                                        >
                                            leave
                                        </Label.Detail>
                                    )}
                                </Label>
                            ))}
                        </div>
                    </Segment>
                </div>
                <div className="col display-flex flex-two">
                    <Chat chat={ctx.chat} sendMessage={ctx.sendMessage} />
                </div>
            </div>
        )
    }

    private leveRoomHandle = (e: MouseEvent) => {
        e.preventDefault()
        this.props.ctx.sendMessage({
            event: "leave_room",
            payload: null,
        })
    }
}

export default Room
