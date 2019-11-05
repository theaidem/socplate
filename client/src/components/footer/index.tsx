import React, { Component } from "react"
import { Label } from "semantic-ui-react"

class Footer extends Component<{ online: number }> {
    public render() {
        const { online } = this.props
        return (
            <footer className="col">
                <Label>
                    online
                    <Label.Detail>{online}</Label.Detail>
                </Label>
            </footer>
        )
    }
}

export default Footer
