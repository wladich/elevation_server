package dem

const HgtSize = 1201
const HgtSplitParts = 4

const TileSize = (HgtSize-1)/HgtSplitParts + 1
const TilePointsN = TileSize * TileSize
const TileBytes = TilePointsN * 2
