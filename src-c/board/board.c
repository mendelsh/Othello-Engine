#include <stdio.h>
#include <stdlib.h>
#include "board.h"

#define SIZE 8
const int shifts[SIZE] = {1, -1, 8, -8, 9, 7, -7, -9};
const bitboard masks[SIZE] = {
	NOT_FILE7,
	NOT_FILE0,
	NOT_RANK7,
	NOT_RANK0,
	NOT_FILE7 & NOT_RANK7,
	NOT_FILE0 & NOT_RANK7,
	NOT_FILE7 & NOT_RANK0,
	NOT_FILE0 & NOT_RANK0
};

#define BLACK_START ( ((uint64_t)1 << (3 + 8*4)) | ((uint64_t)1 << (4 + 8*3)) )
#define WHITE_START ( ((uint64_t)1 << (3 + 8*3)) | ((uint64_t)1 << (4 + 8*4)) )

Board board_create() {
    return (Board){
        .black = BLACK_START,
        .white = WHITE_START,
        .turn = BLACK_TURN
    };
}

void board_display(Board* b) {
	for (int row = 7; row >= 0; row--) {
		printf("%d ", row);
		for (int col = 7; col >= 0; col--) {
			int idx = row*8 + col;
			bitboard mask = (uint64_t)1 << idx;

			if ((b->black & mask) != 0)
				printf("○ ");
			else if ((b->white & mask) != 0)
				printf("● ");
			else
				printf(". ");
		}
		printf("\n");
	}
	printf("  7 6 5 4 3 2 1 0");
	printf("\n");
}

void bit_to_square(uint64_t move, int* x, int* y) {
	int idx = __builtin_ctzll(move);
    *x = idx % 8;
    *y = idx / 8;
}

void count_discs(Board* b, int* black, int* white) {
    *black = __builtin_popcountll(b->black);
    *white = __builtin_popcountll(b->white);
}

#define REPEAT_6(x) x x x x x x

bitboard moves(bitboard player, bitboard opp) {
	bitboard empty = ~(player | opp);
	bitboard moves = 0;

	for (int i = 0; i < SIZE; i++) {
		const bitboard mask = masks[i];
		const int shift = shifts[i];
		bitboard positions = player & mask;
		bitboard discs = 0;

		if (shift > 0) {
			discs = (positions << shift) & opp;

			REPEAT_6({
				positions = discs & mask;
				discs |= (positions << shift) & opp;
			})

			moves |= ((discs & mask) << shift) & empty;
		} else {
			const int s = -shift;
			discs = (positions >> s) & opp;

			REPEAT_6({
				positions = discs & mask;
				discs |= (positions >> s) & opp;
			})

			moves |= ((discs & mask) >> s) & empty;
		}
	}	
	return moves;
}

bitboard flip(bitboard player, bitboard opp, bitboard move) {
	bitboard to_flip = 0;

	for (int i = 0; i < SIZE; i++) {
		const bitboard mask = masks[i];
		const int shift = shifts[i];

		bitboard flips = 0;
		bitboard ray = shift > 0 ? (move & mask) << shift : (move & mask) >> -shift;

		while ((ray & opp) != 0) {
			flips |= ray;
			ray = shift > 0 ? (ray & mask) << shift : (ray & mask) >> -shift;
		}

		if ((ray & player) != 0)
			to_flip |= flips;
	}
	return to_flip;
}

bool play_xy(Board* b, int x, int y) {
	bitboard move = square_to_bit(x, y);

	if ((move & (b->black | b->white)) != 0) return false;
	
	bitboard* player = b->turn == BLACK_TURN ? &b->black : &b->white;
	bitboard* opp = b->turn == WHITE_TURN ? &b->black : &b->white;

	bitboard flipped = flip(*player, *opp, move);

	if (flipped == 0) return false;

	*player |= move | flipped;
	*opp &= ~flipped;

	if (moves(*opp, *player) == 0) return true;	

	flip_turn(b);
	return true;
}

uint64_t count_positions_simple(bitboard player, bitboard opp, unsigned int depth) {
	if (depth == 0 || game_over(player, opp)) return 1;

	bitboard legal_moves = moves(player, opp);

	if (legal_moves == 0) {
		if (moves(opp, player) == 0) return 1;
		return count_positions_simple(opp, player, depth);
	}

	uint64_t count = 0;

	while (legal_moves != 0) {
		bitboard move = (bitboard)1 << __builtin_ctzll(legal_moves);
		legal_moves &= ~move;

		bitboard flipped = flip(player, opp, move);
		bitboard new_player = player | move | flipped;
		bitboard new_opp = opp & ~flipped;

		count += count_positions_simple(new_opp, new_player, depth-1);
	}
	return count;
}

bool play_random_move(Board* b) {
	bitboard legal_moves = b->turn == BLACK_TURN ? moves(b->black, b->white) : moves(b->white, b->black);

	if (legal_moves == 0) return false;

	int count = __builtin_popcountll(legal_moves);
	int choice = rand() % count;

    for (int i = 0; i < 64; i++) {
        if ((legal_moves >> i) & 1) {
            if (choice == 0) return play_xy(b, i%8, i/8);
            
            choice--;
        }
    }

    return false;
}