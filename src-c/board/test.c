#include <stdio.h>
#include <stdlib.h>
#include <time.h>
#include "board.h"

int main() {
    // Seed the random number generator
    srand((unsigned int)time(NULL));

    // Create initial board
    Board b = board_create();

    // Display initial state
    printf("Initial board:\n");
    board_display(&b);

    // Play up to 60 random moves or until the game ends
    int moves_played = 0;
    while (!board_game_over(&b) && moves_played < 60) {
        if (!play_random_move(&b)) {
            printf("No valid moves for %s, skipping turn.\n", b.turn == BLACK_TURN ? "black" : "white");
            flip_turn(&b);  // Force skip if no move
            continue;
        }

        printf("After move %d by %s:\n", moves_played + 1, b.turn == WHITE_TURN ? "black" : "white");
        board_display(&b);
        moves_played++;
    }

    // Count pieces
    int black, white;
    count_discs(&b, &black, &white);

    // Show results
    printf("\nGame over!\n");
    printf("Black: %d\n", black);
    printf("White: %d\n", white);
    if (black > white)
        printf("Black wins!\n");
    else if (white > black)
        printf("White wins!\n");
    else
        printf("It's a tie!\n");

    return 0;
}
