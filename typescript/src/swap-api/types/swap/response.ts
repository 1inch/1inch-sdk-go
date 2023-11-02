export interface QuoteResponse {
    fromToken?: TokenInfo
    toToken?: TokenInfo
    toAmount: string
    protocols?: SelectedProtocol[][]
    gas?: number
}

export interface TokenInfo {
    /** @example 0x111111111117dc0aa78b770fa6a738034120c302 */
    address: string
    /** @example 1INCH */
    symbol: string
    /** @example 1INCH Token */
    name: string
    /** @example 18 */
    decimals: number
    /** @example https://tokens.1inch.io/0x111111111117dc0aa78b770fa6a738034120c302.png */
    logoURI: string
    domainVersion?: string
    /** @example true */
    eip2612?: boolean
    /** @example false */
    isFoT?: boolean
    /**
     * @example [
     *   "tokens"
     * ]
     */
    tags?: string[]
}

export interface SelectedProtocol {
    /**
     * @description Protocol id
     * @example UNISWAP_V3
     */
    name: string
    /**
     * @description Protocol share
     * @example 100
     */
    part: number
    /**
     * @description Source token to swap on protocol
     * @example 0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee
     */
    fromTokenAddress: string
    /**
     * @description Destination token to swap on protocol
     * @example 0x111111111117dc0aa78b770fa6a738034120c302
     */
    toTokenAddress: string
}
