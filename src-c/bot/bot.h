#ifndef BOT_H
#define BOT_H



#include "../board/board.h"
#include <stdlib.h>


#define MAX_MOVES 60

#define MAKE_ALPHA_BETA_FUNC(FUNC_NAME, EVAL_FUNC_NAME)                                             \
            int FUNC_NAME(bitboard player, bitboard opp, int depth, int alpha, int beta) {          \
                if (depth == 0 || game_over(player, opp)) {                                         \
                    return EVAL_FUNC_NAME(player, opp);                                             \
                }                                                                                   \
                                                                                                    \
                bitboard legal_moves = moves(player, opp);                                          \
                if (legal_moves == 0) {                                                             \
                    if (moves(opp, player) == 0) {                                                  \
                        return EVAL_FUNC_NAME(player, opp);                                         \
                    }                                                                               \
                    return -FUNC_NAME(opp, player, depth, -beta, -alpha);                           \
                }                                                                                   \
                                                                                                    \
                while (legal_moves != 0) {                                                          \
                    bitboard move = (bitboard)1 << __builtin_ctzll(legal_moves);                    \
                    legal_moves &= ~move;                                                           \
                                                                                                    \
                    bitboard flips = flip(player, opp, move);                                       \
                    bitboard new_player = player | move | flips;                                    \
                    bitboard new_opp = opp & ~flips;                                                \
                                                                                                    \
                    int score = -FUNC_NAME(new_opp, new_player, depth - 1, -beta, -alpha);          \
                    if (score > alpha) {                                                            \
                        alpha = score;                                                              \
                    }                                                                               \
                    if (alpha >= beta) {                                                            \
                        break;                                                                      \
                    }                                                                               \
                }                                                                                   \
                return alpha;                                                                       \
            }

#define BLACK true
#define WHITE false

#define MAX_EVAL 100000
#define MIN_EVAL -MAX_EVAL


typedef bool Side;

typedef struct
{
    int x;
    int y;
    int score;
} Move;

#define EMPTY_MOVE ((Move){-1, -1, 0})

typedef struct
{
    Move move[MAX_MOVES];
    int len;
} MoveSlice;

MoveSlice* MoveSlice_new(void);
void MoveSlice_safe_free(MoveSlice** slice_ptr);
void MoveSlice_add(MoveSlice* slice, Move move);
Move MoveSlice_get(MoveSlice* slice, int idx);
void MoveSlice_clear(MoveSlice* slice);

MoveSlice* search(bitboard player,
                  bitboard opp,
                  int depth,
                  int (*eval)(bitboard player, bitboard opp, int depth));

typedef struct
{
    int (*eval)(bitboard player, bitboard opp, int depth);
    MoveSlice* (*search)(bitboard player,
                         bitboard opp,
                         int depth,
                         int (*eval)(bitboard player, bitboard opp, int depth));
} IntEvaluation;

Move select_move(IntEvaluation* int_evaluation, bitboard player, bitboard opp, int depth);

typedef struct
{
    bitboard* player;
    bitboard* opponent;
    Side side;
    IntEvaluation* int_evaluation;
} Bot;

Bot Bot_create(Board* board, Side side, IntEvaluation* int_evaluation);
Move Bot_select_move(Bot* bot, int depth);

int Pengwin_static_eval(bitboard player, bitboard opp);
int Pengwin_eval(bitboard player, bitboard opp, int depth);
Bot Pengwin_create(Board* board, Side side);

int Greedy_static_eval(bitboard player, bitboard opp);
int Greedy_eval(bitboard player, bitboard opp, int depth);
Bot Greedy_create(Board* board, Side side);



#endif // BOT_H