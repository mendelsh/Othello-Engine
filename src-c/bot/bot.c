#include "bot.h"
#include "stable.h"
#include <stdio.h>

MoveSlice* MoveSlice_new() {
    MoveSlice* slice = (MoveSlice*)malloc(sizeof(MoveSlice));
    if (slice == NULL) {
        return NULL;
    }
    slice->len = 0;
    return slice;
}

void MoveSlice_safe_free(MoveSlice** slice_ptr) {
    if (slice_ptr != NULL && *slice_ptr != NULL) {
        free(*slice_ptr);
        *slice_ptr = NULL;
    }
}

void MoveSlice_add(MoveSlice* slice, Move move) {
    slice->move[slice->len] = move;
    slice->len++;
}

Move MoveSlice_get(MoveSlice* slice, int idx) {
    return slice->move[idx];
}

void MoveSlice_clear(MoveSlice* slice) {
    slice->len = 0;
}

int compare(const void *a, const void *b) {
    return ((Move *)b)->score - ((Move *)a)->score;
}

MoveSlice* search(bitboard player,
                  bitboard opp,
                  int depth,
                  int (*eval)(bitboard player, bitboard opp, int depth))
{
    MoveSlice* slice = MoveSlice_new();
    if (slice == NULL) return NULL;
    
    bitboard legal_moves = moves(player, opp);
    while (legal_moves != 0) {
        bitboard move = (bitboard)1 << __builtin_ctzll(legal_moves);
        legal_moves &= ~move;
        
        bitboard flipped = flip(player, opp, move);
        bitboard new_player = player | move | flipped;
        bitboard new_opp = opp & ~flipped;
        
        int score = -eval(new_opp, new_player, depth - 1);
        int x, y;
        bit_to_square(move, &x, &y);
        MoveSlice_add(slice, (Move){x, y, score});
    }
    qsort(slice->move, slice->len, sizeof(Move), compare);

    return slice;
}

Move select_move(IntEvaluation* int_evaluation, bitboard player, bitboard opp, int depth) {
    MoveSlice* slice = int_evaluation->search(player, opp, depth, int_evaluation->eval);
    if (slice == NULL) return EMPTY_MOVE;

    int best_score = slice->move[0].score;
    int i = 0;
    for (i = 1; i < slice->len; i++) {
        if (slice->move[i].score < best_score)
            break;
    }

    Move move = MoveSlice_get(slice, rand() % i);
    MoveSlice_safe_free(&slice);
    return move;
}

Bot Bot_create(Board* board, Side side, IntEvaluation* int_evaluation) {
    Bot bot;
    bot.side = side;
    bot.player = side == BLACK ? &board->black : &board->white;
    bot.opponent = side == WHITE ? &board->black : &board->white;
    bot.int_evaluation = int_evaluation;
    return bot;
}

Move Bot_select_move(Bot* bot, int depth) {
    return select_move(bot->int_evaluation, *bot->player, *bot->opponent, depth);
}


// Pengwin Bot
int Pengwin_static_eval(bitboard player, bitboard opp) {
    if (game_over(player, opp)) {
        return 1000 * (__builtin_popcountll(player) - __builtin_popcountll(opp));
    }
    int flexibility = __builtin_popcountll(moves(player, opp)) - __builtin_popcountll(moves(opp, player));
    int cp = 0, co = 0;
    stable_discs_count(player, opp, &cp, &co);
    int stability = 20 * (cp - co);
    return flexibility + stability;
}

MAKE_ALPHA_BETA_FUNC(Pengwin_alpha_beta, Pengwin_static_eval)

int Pengwin_eval(bitboard player, bitboard opp, int depth) {
    return Pengwin_alpha_beta(player, opp, depth, MIN_EVAL, MAX_EVAL);
}

Bot Pengwin_create(Board* board, Side side) {
    static IntEvaluation eval = {Pengwin_eval, search};
    return Bot_create(board, side, &eval);
}

// Greedy Bot
int Greedy_static_eval(bitboard player, bitboard opp) {
    int quantity =__builtin_popcountll(player) - __builtin_popcountll(opp);
    if (game_over(player, opp)) {
        return 1000 * quantity;
    }
    int cp = 0, co = 0;
    stable_discs_count(player, opp, &cp, &co);
    int stability = 500 * (cp - co);
    return stability + quantity;
}

MAKE_ALPHA_BETA_FUNC(Greedy_alpha_beta, Greedy_static_eval)

int Greedy_eval(bitboard player, bitboard opp, int depth) {
    return Greedy_alpha_beta(player, opp, depth, MIN_EVAL, MAX_EVAL);
}

Bot Greedy_create(Board* board, Side side) {
    static IntEvaluation eval = {Greedy_eval, search};
    return Bot_create(board, side, &eval);
}