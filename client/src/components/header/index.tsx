import React, { Component } from "react"
import { Label, Header, Icon } from "semantic-ui-react"
import { User } from "../../types/user"

class Head extends Component<{ user: User }> {
    public render() {
        const { user } = this.props
        return (
            <header className="col display-flex">
                <Header as="h2">
                    <Icon name="user secret" />
                    socplate
                </Header>
                <Label as="a" image={true}>
                    <img alt={user.name} src={user.avatar} />
                    {user.name}
                </Label>
            </header>
        )
    }
}

export default Head
