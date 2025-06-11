#ifndef STABLE_H
#define STABLE_H


#include "../board/board.h"


#define CORNERS ((bitboard)0b1000000100000000000000000000000000000000000000000000000010000001)

static inline bitboard stable_discs(bitboard board) {
    return board & CORNERS;
}

static inline void stable_discs_count(bitboard player, bitboard opp, int* player_count, int* opp_count) {
    *player_count = __builtin_popcountll(stable_discs(player));
    *opp_count = __builtin_popcountll(stable_discs(opp));
}



#endif // STABLE_H