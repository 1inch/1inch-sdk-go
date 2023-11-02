import BN from 'bn.js'
import {NATIVE_CURRENCY} from './constants'

export const isNativeCurrency = (address: string): boolean =>
    address.toLowerCase() === NATIVE_CURRENCY

export function toSec(time: number | string | Date): number {
    const t = time instanceof Date ? time.getTime() : time

    return Math.floor(+t / 1000)
}

export function toBN(val: number | string): BN {
    if (typeof val === 'number') {
        if (!Number.isSafeInteger(val)) {
            throw new Error('integer is not safe')
        }

        return new BN(val)
    }

    if (val.startsWith('0x')) {
        return new BN(trim0x(val), 'hex')
    }

    return new BN(val)
}

export function trim0x(data: string): string {
    if (data.startsWith('0x')) {
        return data.substring(2)
    }

    return data
}

export function add0x(data: string): string {
    if (data.includes('0x')) {
        return data
    }

    return '0x' + data
}

export function getCrypto(): Crypto {
    if (typeof window !== 'undefined') {
        return window.crypto
    } else {
        // eslint-disable-next-line @typescript-eslint/no-var-requires
        return require('crypto').webcrypto as unknown as Crypto
    }
}

export function concatQueryParams<
    T extends Record<string | number, string | string[] | number | boolean>
>(params: T): string {
    if (!params) {
        return ''
    }

    const keys = Object.keys(params)

    if (keys.length === 0) {
        return ''
    }

    return (
        '?' +
        keys
            .reduce((a, k) => {
                if (!params[k]) {
                    return a
                }

                const value = params[k]
                a.push(
                    k +
                        '=' +
                        encodeURIComponent(
                            Array.isArray(value) ? value.join(',') : value
                        )
                )

                return a
            }, [] as string[])
            .join('&')
    )
}
