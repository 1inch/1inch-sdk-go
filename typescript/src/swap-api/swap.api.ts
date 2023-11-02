// import {SwapRequest} from './swap.request'
import {SwapApiConfig} from './types'
import {AxiosProviderConnector, HttpProviderConnector} from '../connectors'
import {concatQueryParams} from '../utils'
import {QuoteRequest} from './types/swap/request'

export class SwapApi {
    constructor(
        private readonly config: SwapApiConfig,
        private readonly httpClient: HttpProviderConnector
    ) {}

    static new(
        config: SwapApiConfig,
        httpClient: HttpProviderConnector = new AxiosProviderConnector(
            config.authKey
        )
    ): SwapApi {
        return new SwapApi(config, httpClient)
    }

    quote(params: QuoteRequest): Promise<void> {
        const err = params.validate()

        if (err) {
            throw new Error(err)
        }

        // eslint-disable-next-line @typescript-eslint/ban-ts-comment
        // @ts-ignore
        const queryParams = concatQueryParams(params.build())
        const url = `${this.config.url}/swap/${this.config.version}/${this.config.network}/quote?${queryParams}`

        return this.httpClient.post(url, queryParams)
    }

    swap(): Promise<void> {
        return Promise.resolve()
    }

    spender(): Promise<void> {
        return Promise.resolve()
    }

    transaction(): Promise<void> {
        return Promise.resolve()
    }

    allowance(): Promise<void> {
        return Promise.resolve()
    }

    liquiditySources(): Promise<void> {
        return Promise.resolve()
    }

    tokens(): Promise<void> {
        return Promise.resolve()
    }

    // submit(params: SwapRequest): Promise<void> {
    //     const url = `${this.config.url}/v1.0/${this.config.network}/order/submit`
    //
    //     return this.httpClient.post(url, params)
    // }
    //
    // submitBatch(params: SwapRequest[]): Promise<void> {
    //     const url = `${this.config.url}/v1.0/${this.config.network}/order/submit/many`
    //
    //     return this.httpClient.post(url, params)
    // }
}
