#include <stdio.h>
#include <stdlib.h>
#include <time.h>
#include "bot.h"


#define DEPTH 4

typedef enum 
{
    BLACK_WIN,
    WHITE_WIN,
    DRAW
} Result;

int main() {
    // Seed the random number generator
    srand((unsigned int)time(NULL));

    // Create initial board
    Board b = board_create();

    // Display initial state
    printf("Initial board:\n");
    board_display(&b);

    // Statistics
    Bot black_bot, white_bot;
    int pengwin_wins = 0, greedy_wins = 0, ties = 0;
    int black_wins = 0, white_wins = 0;

    for (int i = 0; i < 1000; i++) {
        Board b = board_create();
        if (i % 2 == 0) {
            black_bot = Pengwin_create(&b, BLACK_TURN);
            white_bot = Greedy_create(&b, WHITE_TURN);
        } else {
            black_bot = Greedy_create(&b, BLACK_TURN);
            white_bot = Pengwin_create(&b, WHITE_TURN);
        }
        while (!board_game_over(&b)) {
            Move m = Bot_select_move(b.turn == BLACK_TURN ? &black_bot : &white_bot, DEPTH);
            play_xy(&b, m.x, m.y);
        }
        int black, white;
        count_discs(&b, &black, &white);
        if (black == white) {
            ties++;
        } else if (black > white) {
            if (i % 2 == 0) {
                pengwin_wins++;
            } else {
                greedy_wins++;
            }
            black_wins++;
        } else {
            if (i % 2 == 0) {
                greedy_wins++;

            } else {
                pengwin_wins++;
            }
            white_wins++;
        }
    }

    printf("Pengwin wins: %d (%.2f%%)\n", pengwin_wins, pengwin_wins / 10.0);
    printf("Greedy wins: %d (%.2f%%)\n", greedy_wins, greedy_wins / 10.0);
    printf("Ties: %d (%.2f%%)\n", ties, ties / 10.0);
    printf("Black wins: %d (%.2f%%)\n", black_wins, black_wins / 10.0);
    printf("White wins: %d (%.2f%%)\n", white_wins, white_wins / 10.0);

/*
    Bot pengwin = Pengwin_create(&b, BLACK_TURN);
    Bot greedy = Greedy_create(&b, WHITE_TURN);

    while (!board_game_over(&b)) {
        Move m = b.turn == BLACK_TURN ? Bot_select_move(&pengwin, DEPTH) : Bot_select_move(&greedy, DEPTH);
        play_xy(&b, m.x, m.y);
        printf("%s turn:\n", b.turn == BLACK_TURN ? "black(pengwin)" : "white(greedy)");
        board_display(&b);
    }
    
    printf("Game over!\n");
    printf("Final board:\n");
    board_display(&b);
    
    int black, white;
    count_discs(&b, &black, &white);
    printf("Black: %d\n", black);
    printf("White: %d\n", white);

    Result result;
    if (black == white) {
        result = DRAW;
    } else if (black > white) {
        result = BLACK_WIN;
    } else {
        result = WHITE_WIN;
    }

    switch (result) {
        case BLACK_WIN:
            printf("Black (%s) wins!\n", pengwin.side == BLACK ? "pengwin" : "greedy");
            break;
        case WHITE_WIN:
            printf("White (%s) wins!\n", greedy.side == WHITE ? "greedy" : "pengwin");
            break;
        case DRAW:
            printf("Draw!\n");
            break;
    }
*/


    return 0;
}
