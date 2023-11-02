import {NATIVE_CURRENCY, ZERO_ADDRESS} from '../../../constants'
import {
    hasDuplicates,
    isValidAddress,
    isValidAmount
} from '../../../validations'

export class QuoteRequest {
    public readonly src: string

    public readonly dst: string

    public readonly amount: string

    public readonly fee?: number

    public readonly gasPrice?: string

    public readonly complexityLevel?: number

    public readonly parts?: number

    public readonly mainRouteParts?: number

    public readonly gasLimit?: number

    public readonly connectorTokens?: string[]

    constructor(params: QuoterRequestParamsRaw) {
        this.src = params.src
        this.dst = params.dst
        this.amount = params.amount

        if (params.fee) {
            this.fee = params.fee
        }

        if (params.gasPrice) {
            this.gasPrice = params.gasPrice
        }

        if (params.complexityLevel) {
            this.complexityLevel = params.complexityLevel
        }

        if (params.parts) {
            this.parts = params.parts
        }

        if (params.mainRouteParts) {
            this.mainRouteParts = params.mainRouteParts
        }

        if (params.gasLimit) {
            this.gasLimit = params.gasLimit
        }

        if (params.gasLimit) {
            this.gasLimit = params.gasLimit
        }

        if (params.connectorTokens) {
            this.connectorTokens = params.connectorTokens
        }
    }

    static new(params: QuoterRequestParamsRaw): QuoteRequest {
        return new QuoteRequest(params)
    }

    validate(): string | null {
        if (this.src === ZERO_ADDRESS || this.dst === ZERO_ADDRESS) {
            return `replace ${ZERO_ADDRESS} with ${NATIVE_CURRENCY}`
        }

        if (this.src === this.dst) {
            return 'fromTokenAddress and toTokenAddress should be different'
        }

        if (!isValidAddress(this.src)) {
            return `${this.src} is invalid fromTokenAddress`
        }

        if (!isValidAddress(this.dst)) {
            return `${this.dst} is invalid toTokenAddress`
        }

        if (this.fee && (this.fee > 3 || this.fee < 0)) {
            return `${this.fee} is an invalid fee, max = 3 and min = 0`
        }

        if (this.gasPrice && !isValidAmount(this.gasPrice)) {
            return `${this.gasPrice} is invalid gasPrice`
        }

        if (
            this.complexityLevel &&
            (this.complexityLevel > 2 || this.complexityLevel < 0)
        ) {
            return `${this.complexityLevel} is invalid, max = 2 and min = 0 `
        }

        if (this.parts && (this.parts > 50 || this.parts < 0)) {
            return `${this.parts} is invalid, max = 50 and min = 0 `
        }

        if (
            this.mainRouteParts &&
            (this.mainRouteParts > 10 || this.mainRouteParts < 0)
        ) {
            return `${this.mainRouteParts} is invalid, max = 10 and min = 0 `
        }

        if (
            this.gasLimit &&
            (this.gasLimit < 100000 || this.gasLimit > 11500000)
        ) {
            return `${this.gasLimit} is invalid, max = 11500000 and min = 100000 `
        }

        if (this.connectorTokens) {
            if (this.connectorTokens.length > 5) {
                return `${this.connectorTokens.length} length is invalid, max = 5`
            }

            if (hasDuplicates(this.connectorTokens)) {
                return `connectorTokens has duplicates`
            }

            for (const connectorToken of this.connectorTokens) {
                if (!isValidAddress(connectorToken)) {
                    return `${connectorToken} is invalid connectorToken`
                }
            }
        }

        return null
    }

    build(): QuoterRequestParams {
        const params: QuoterRequestParams = {
            src: this.src,
            dst: this.dst,
            amount: this.amount
        }

        if (this.fee) {
            params.fee = this.fee
        }

        if (this.gasPrice) {
            params.gasPrice = this.gasPrice
        }

        if (this.complexityLevel) {
            params.complexityLevel = this.complexityLevel
        }

        if (this.parts) {
            params.parts = this.parts
        }

        if (this.mainRouteParts) {
            params.mainRouteParts = this.mainRouteParts
        }

        if (this.gasLimit) {
            params.gasLimit = this.gasLimit
        }

        if (this.connectorTokens) {
            params.connectorTokens = this.connectorTokens.join(',')
        }

        return params
    }
}

export interface QuoterRequestParamsRaw {
    /** @example 0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee */
    src: string
    /** @example 0x111111111117dc0aa78b770fa6a738034120c302 */
    dst: string
    /** @example 10000000000000000 */
    amount: string
    /** @description All supported liquidity sources by default */
    protocols?: string
    /** @description Partner fee. min: 0; max: 3 Should be the same for /quote and /swap */
    fee?: number
    /** @description Network price per gas in wei. By default fast network gas price */
    gasPrice?: string
    complexityLevel?: number
    parts?: number
    mainRouteParts?: number
    gasLimit?: number
    /** @description Return fromToken and toToken info in response */
    includeTokensInfo?: boolean
    /** @description Return used swap protocols in response */
    includeProtocols?: boolean
    /** @description Return approximated gas in response */
    includeGas?: boolean
    connectorTokens?: string[]
}

export interface QuoterRequestParams {
    /** @example 0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee */
    src: string
    /** @example 0x111111111117dc0aa78b770fa6a738034120c302 */
    dst: string
    /** @example 10000000000000000 */
    amount: string
    /** @description All supported liquidity sources by default */
    protocols?: string
    /** @description Partner fee. min: 0; max: 3 Should be the same for /quote and /swap */
    fee?: number
    /** @description Network price per gas in wei. By default fast network gas price */
    gasPrice?: string
    complexityLevel?: number
    parts?: number
    mainRouteParts?: number
    gasLimit?: number
    /** @description Return fromToken and toToken info in response */
    includeTokensInfo?: boolean
    /** @description Return used swap protocols in response */
    includeProtocols?: boolean
    /** @description Return approximated gas in response */
    includeGas?: boolean
    connectorTokens?: string
}
