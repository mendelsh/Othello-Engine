#ifndef BOARD_H
#define BOARD_H



#include <stdint.h>
#include <stdbool.h>

#define FILE0 (uint64_t)0x0101010101010101
#define FILE7 (uint64_t)0x8080808080808080
#define RANK0 (uint64_t)0x00000000000000FF
#define RANK7 (uint64_t)0xFF00000000000000

#define NOT_FILE0 (~FILE0 & (uint64_t)0xFFFFFFFFFFFFFFFF)
#define NOT_FILE7 (~FILE7 & (uint64_t)0xFFFFFFFFFFFFFFFF)
#define NOT_RANK0 (~RANK0 & (uint64_t)0xFFFFFFFFFFFFFFFF)
#define NOT_RANK7 (~RANK7 & (uint64_t)0xFFFFFFFFFFFFFFFF)

#define BLACK_TURN true
#define WHITE_TURN false

typedef uint64_t bitboard;
typedef bool Turn;

typedef struct
{
	bitboard black;
	bitboard white;
	Turn turn;
} Board;

Board board_create();
void board_display(Board*);

static inline bitboard empty(Board* b) {
    return ~(b->black | b->white);
}

static inline bitboard square_to_bit(int x, int y) {
	return (uint64_t)1 << (y * 8 + x);
}

void bit_to_square(uint64_t move, int* x, int* y);

void count_discs(Board* b, int* black, int* white);

bitboard moves(bitboard player, bitboard opp);
bitboard flip(bitboard player, bitboard opp, bitboard move);

static inline void flip_turn(Board* b) {
	b->turn = !b->turn;
}

bool play_xy(Board* b, int x, int y);

static inline bool game_over(bitboard player1, bitboard player2) {
	return moves(player1, player2) == 0 && moves(player2, player1) == 0;
}
static inline bool board_game_over(Board* b) {
	return game_over(b->black, b->white);
}

uint64_t count_positions_simple(bitboard player, bitboard opp, unsigned int depth);
bool play_random_move(Board*);



#endif // BOARD_H