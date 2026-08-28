package orderbook

// Compile-time assertion that the deprecated GetSaltParams name still resolves.
var _ GetSaltParams = GenerateSaltWithFeesParams{}
