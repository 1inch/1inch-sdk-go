import {NetworkEnum} from '../constants';

export type SwapRequestParams = {
    signature: string
    quoteId: string
}

export type SwapApiConfig = {
    network: NetworkEnum
    url: string
    authKey?: string
}
