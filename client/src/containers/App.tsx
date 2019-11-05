import React from "react"
import { CtxProps } from "../types/context"
import Config from "../services/Config"
import withCtx from "../contexts/App"
import Layout from "../components/layout"
import Header from "../components/header"
import Footer from "../components/footer"
import Lobby from "../components/lobby"
import Room from "../components/room"

class App extends React.Component<CtxProps> {
    public componentDidMount() {}

    public componentWillUpdate(nextProps: CtxProps) {
        // tslint:disable-next-line: no-console
        if (!Config.isProd()) console.debug("context", nextProps)
    }

    public render() {
        const { ctx } = this.props
        if (!ctx.user) return null
        return (
            <Layout ctx={ctx}>
                <div id={`user_${ctx.user.id}`} className="app display-flex">
                    <div className="wrapper display-flex flex-one flex-direction-col">
                        <Header user={ctx.user} />
                        <main className="content flex-one display-flex">
                            {this.renderStage()}
                        </main>
                        <Footer online={ctx.online} />
                    </div>
                </div>
            </Layout>
        )
    }

    private renderStage = () => {
        const { ctx } = this.props
        switch (ctx.stage) {
            case "lobby":
                return <Lobby ctx={ctx} />
            case "room":
                return <Room ctx={ctx} />
            default:
                return null
        }
    }
}

export default withCtx(App)
