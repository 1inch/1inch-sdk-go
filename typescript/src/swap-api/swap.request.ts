import {SwapRequestParams} from './types'
import {LimitOrderV3Struct} from '../../limit-order'

export class SwapRequest {
    public readonly order: LimitOrderV3Struct

    public readonly signature: string

    public readonly quoteId: string

    constructor(params: SwapRequestParams) {
        this.order = params.order
        this.signature = params.signature
        this.quoteId = params.quoteId
    }

    static new(params: SwapRequestParams): SwapRequest {
        return new SwapRequest(params)
    }

    build(): SwapRequestParams {
        return {
            order: this.order,
            signature: this.signature,
            quoteId: this.quoteId
        }
    }
}
