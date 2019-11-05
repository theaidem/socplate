import React, { Component, createContext } from "react"
import { encode, decode } from "@msgpack/msgpack"

import { CtxState, CtxActions, DimmerStage } from "../types/context"
import { IncomingMessage, OutcomingMessage } from "../types/message"
import API from "../services/API"
import Config from "../services/Config"

const defaultState: CtxState = {
    isInitialized: false,
    dimming: {
        active: true,
        stage: "load",
    },
    user: null,
    stage: "lobby",
    room: null,
    rooms: [],
    chat: [],
    online: 0,
}

const appCtx = createContext(defaultState)

export class AppStore extends Component {
    public state = defaultState
    public socket?: WebSocket = undefined

    public componentDidMount = async () => {
        await this.init()
    }

    public render = () => {
        const { children } = this.props
        const { isInitialized, user, stage } = this.state
        if (stage === "error") return "error"

        return (
            <appCtx.Provider
                value={{
                    ...this.exportState(),
                    ...this.exportActions(),
                }}
            >
                {isInitialized && user ? children : "loading"}
            </appCtx.Provider>
        )
    }

    private exportState = (): CtxState => ({ ...this.state })

    private exportActions = (): CtxActions => ({
        sendMessage: this.sendMessage,
    })

    private toggleDimming = (
        active: boolean,
        stage: DimmerStage = null,
        timeout: number = 0,
    ) => {
        setTimeout(() => {
            this.setState({ dimming: { active, stage } })
        }, timeout)
    }

    private restartApp = (timeout = 3000) => {
        setTimeout(() => this.init(), timeout)
    }

    private sendMessage = (msg: OutcomingMessage): void => {
        if (!this.socket) return
        this.socket.send(encode(msg))
    }

    private init = async () => {
        const resp = await API.get("/index")
        if (!resp) return this.setState({ stage: "error" }, this.restartApp)

        this.socket = new WebSocket(
            `${Config.websocketAddr()}?ticket=${resp.data.ticket}`,
        )
        this.socket.binaryType = "arraybuffer"
        this.socket.onopen = this.socketOnOpen
        this.socket.onclose = this.socketOnClose
        this.socket.onerror = this.socketOnError
        this.socket.onmessage = this.socketOnMessage
    }

    private socketOnOpen = (evt: Event): void => {
        this.setState({ isInitialized: true, stage: "lobby" }, () =>
            this.toggleDimming(false, null, 1500),
        )
    }

    private socketOnClose = (evt: Event): void => {
        this.toggleDimming(true, "error")
        this.restartApp()
    }

    private socketOnError = (evt: Event): void => {
        // tslint:disable-next-line: no-console
        console.log("error", evt)
        this.toggleDimming(true, "error")
        this.restartApp()
    }

    private socketOnMessage = (evt: MessageEvent) => {
        const binData = new Uint8Array(evt.data)
        const message = decode(binData)
        this.handleMessage(message as IncomingMessage)
    }

    private handleMessage = (message: IncomingMessage): void => {
        // tslint:disable-next-line: no-console
        if (!Config.isProd()) console.debug("message", message)

        const { stage } = this.state
        switch (message.event) {
            case "online":
                this.setState({ online: message.payload })
                break

            case "user":
                this.setState({
                    user: message.payload,
                    stage:
                        !Boolean(message.payload.room) && stage === "room"
                            ? "lobby"
                            : stage,
                })
                break

            case "room_list":
                this.setState({ rooms: message.payload })
                break

            case "start_room":
                this.setState({
                    stage: "room",
                    room: message.payload,
                    rooms: defaultState.rooms,
                    chat: defaultState.chat,
                })
                break

            case "update_room":
                this.setState({ room: message.payload })
                break

            case "chat_messages":
                this.setState({ chat: message.payload.reverse() })
                break

            default:
                // tslint:disable-next-line: no-console
                console.log("Unknown message event", message)
        }
    }
}

// tslint:disable-next-line: variable-name
const withCtx = (Comp: any) => (props: any) => (
    <appCtx.Consumer>{data => <Comp {...props} ctx={data} />}</appCtx.Consumer>
)

export default withCtx
