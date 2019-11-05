import React, { ChangeEvent, Component } from "react"
import { formatDistance } from "date-fns"
import { Header, Segment, Comment, Form, Input, Icon } from "semantic-ui-react"
import { Chat as TChat, ChatMessage } from "../../types/chat"
import { CtxActions } from "../../types/context"

export default class Chat extends Component<{
    chat: TChat
    sendMessage: CtxActions["sendMessage"]
}> {
    public input: Input | null = null
    public stateUpdater: NodeJS.Timeout | null = null
    public state = {
        message: "",
    }

    public componentDidMount() {
        this.stateUpdater = setInterval(() => {
            this.forceUpdate()
        }, 1000)
    }

    public componentWillUnmount() {
        if (this.stateUpdater !== null) clearInterval(this.stateUpdater)
        this.stateUpdater = null
    }

    public render() {
        const { chat } = this.props

        return (
            <Segment className="display-flex flex-one flex-direction-col">
                <Comment.Group
                    className="display-flex flex-one flex-direction-col"
                    minimal={true}
                    size="small"
                >
                    <Header as="h3" dividing={true} size="small">
                        <Icon name="conversation" />
                        Chat
                    </Header>
                    <div className="msg-list display-flex flex-one">
                        {chat.map((msg: ChatMessage) => (
                            <Comment key={msg.created}>
                                <Comment.Avatar as="a" src={msg.from.avatar} />
                                <Comment.Content>
                                    <Comment.Author as="a">
                                        {msg.from.name}
                                    </Comment.Author>
                                    <Comment.Metadata>
                                        <span>
                                            {formatDistance(
                                                Math.floor(
                                                    msg.created / 1000000,
                                                ),
                                                new Date(),
                                                { addSuffix: true },
                                            )}
                                        </span>
                                    </Comment.Metadata>
                                    <Comment.Text>{msg.message}</Comment.Text>
                                </Comment.Content>
                            </Comment>
                        ))}
                    </div>
                    <Form reply={true}>
                        <Input
                            fluid={true}
                            ref={i => (this.input = i)}
                            size="small"
                            value={this.state.message}
                            onChange={this.writeMessageHandle}
                            action={{
                                onClick: this.sendMessageHandle,
                                labelPosition: "right",
                                icon: "send",
                                content: "Send",
                            }}
                            focus={true}
                            placeholder="Write a message..."
                        />
                    </Form>
                </Comment.Group>
            </Segment>
        )
    }

    private writeMessageHandle = (e: ChangeEvent, data: any) => {
        e.preventDefault()
        if (this.input === null) return
        this.setState({ message: data.value }, () => {
            if (this.input !== null) {
                this.input.focus() // ??
            }
        })
    }

    private sendMessageHandle = (e: ChangeEvent, data: object) => {
        const { message } = this.state
        const { sendMessage } = this.props
        e.preventDefault()
        if (message === "") return
        sendMessage({
            event: "chat_message",
            payload: message,
        })
        this.setState({ message: "" })
    }
}
