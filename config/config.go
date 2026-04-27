package config

// Defining size and shape of board
// Stored in Uints right now to make Uint -> Vec easier
// BUT it may be benificial later to store them as ints
// And to make Vec -> Uint easier
const BoardSize uint32 = 512
const LayerSize uint32 = 64
const LineSize uint32 = 8
const SpaceSize uint32 = 1
