import React, { Component } from "react"
import { Dimmer, Loader } from "semantic-ui-react"
import { CtxProps } from "../../types/context"

class Layout extends Component<CtxProps> {
    public render() {
        const { children } = this.props
        const { active } = this.props.ctx.dimming
        return (
            <Dimmer.Dimmable as="div" blurring={true} dimmed={active}>
                <Dimmer active={active} inverted={true}>
                    {this.renderDimmerStage()}
                </Dimmer>
                {children}
            </Dimmer.Dimmable>
        )
    }

    private renderDimmerStage = () => {
        const { stage } = this.props.ctx.dimming
        switch (stage) {
            case "load":
                return <Loader>loading</Loader>
            case "error":
                return <Loader>error</Loader>

            default:
                return null
        }
    }
}

export default Layout
