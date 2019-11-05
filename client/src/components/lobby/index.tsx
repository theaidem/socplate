import React, { Component } from "react"
import { CtxProps } from "../../types/context"
import RoomList from "../room/list"
import Chat from "../chat"

class Lobby extends Component<CtxProps> {
    public render() {
        const { ctx } = this.props
        return (
            <div className="display-flex flex-one">
                <div className="col display-flex flex-one">
                    <RoomList ctx={ctx} />
                </div>
                <div className="col display-flex flex-one">
                    <Chat chat={ctx.chat} sendMessage={ctx.sendMessage} />
                </div>
            </div>
        )
    }
}

export default Lobby
