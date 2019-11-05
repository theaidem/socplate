import axios, { AxiosInstance, AxiosResponse, AxiosRequestConfig } from "axios"
import Config from './Config'

export class API {
    public service: AxiosInstance

    constructor(baseURL?: string) {
        const config = {
            baseURL,
            headers: {},
            timeout: 7000,
        }
        const service = axios.create(config)
        service.interceptors.response.use(this.handleSuccess, this.handleError)
        this.service = service
    }

    public handleSuccess = (response: AxiosResponse): AxiosResponse => response

    public handleError = (error: any) => error.response

    public async get(path: string, headers: any = null) {
        const data = await this.service.get(path)
        return data
    }

    public async patch(
        path: string,
        data: any = null,
        headers: any = null,
    ) {
        const requestConfig: AxiosRequestConfig = {
            method: "PATCH",
            responseType: "json",
            url: path,
            data, headers,
        }

        const response = await this.service.request(requestConfig)
        return response
    }

    public async post(
        path: string | undefined,
        data: any = null,
        headers: any = null,
    ) {
        const requestConfig: AxiosRequestConfig = {
            method: "POST",
            responseType: "json",
            url: path,
            data, headers,
        }

        const response = await this.service.request(requestConfig)
        return response
    }
}

const instance = new API(Config.httpAddr())
export default instance
