package fusionorder

import "math/big"

// Uint16Max is the maximum value for a uint16 (2^16 - 1)
var Uint16Max = new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 16), big.NewInt(1))

// Uint24Max is the maximum value for a uint24 (2^24 - 1)
const Uint24Max = (1 << 24) - 1

// Uint32Max is the maximum value for a uint32 (2^32 - 1)
const Uint32Max = (1 << 32) - 1

// Uint40Max is the maximum value for a uint40 (2^40 - 1) as *big.Int
var Uint40Max = new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 40), big.NewInt(1))
