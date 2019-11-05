
export class Config {

    public isProd = () => process.env.NODE_ENV === "production"

    public httpAddr = () => process.env.REACT_APP_SERVER_HTTP_ADDR

    public websocketAddr = () => process.env.REACT_APP_SERVER_WS_ADDR

}

export default new Config()
