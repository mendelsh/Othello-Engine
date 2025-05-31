
class Board {
    constructor(containerId) {
      this.container = document.getElementById(containerId);
      this.cells = [];
      this.countBlack = 2;
      this.countWhite = 2;
      this.createBoard();
      this.init();
    }
    
    createBoard() {
      for (let i = 63; i >= 0; i--) {
        const cell = document.createElement('div');
        cell.classList.add('cell');
        cell.dataset.index = i;
        cell.textContent = '';
        
        cell.addEventListener('click', () => {
            handleMove(i);
            // console.log(i)
        });
        
        this.container.appendChild(cell);
        this.cells[i] = cell;
      }
    }
    
    updateBoard(blackBoard, whiteBoard) {
        this.countBlack = 0;
        this.countWhite = 0;

        for (let i = 0; i < 64; i++) {
            // console.log(blackBoard.toString(2)); 
            const blackBit = (blackBoard >> BigInt(i)) & 1n;
            const whiteBit = (whiteBoard >> BigInt(i)) & 1n;

            // console.log(`Bit ${i}:`, blackBit.toString());
            
        
            // Clear previous disc if any
            this.cells[i].innerHTML = '';
        
            if (blackBit === 1n) {
                this.countBlack++;
                const disc = document.createElement('div');
                disc.classList.add('disc', 'black');
                this.cells[i].appendChild(disc);
            } else if (whiteBit === 1n) {
                this.countWhite++;
                const disc = document.createElement('div');
                disc.classList.add('disc', 'white');
                this.cells[i].appendChild(disc);
            }
        }
    }      

    async init() {
        try {
            const res = await fetch("/state");
            const data = await res.json();
            blackTurn = data.black_turn;

            const blackBoard = BigInt(data.black);
            const whiteBoard = BigInt(data.white);

            this.updateBoard(blackBoard, whiteBoard);
        } 
        catch (err) {
            console.error("Init error:", err);
        }
      }
}

const board = new Board("board");
const status = document.getElementById("status");
const count = document.getElementById("count");

let blackTurn = true;

function shiftChar(c, n) {
    return String.fromCharCode(c.charCodeAt(0) - n);
}

function indexToSquare(index) {
    let num = Math.floor((index / 8) + 1);
    let char = shiftChar('h', index % 8);
    return `${char}${num}`;
}

async function handleMove(index) {
    const move = indexToSquare(index);
    // console.log(move)
    try {
        const response = await fetch("/move", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ move }),
        });

        if (!response.ok) {
            return;
        }

        const data = await response.json();

        const blackBoard = BigInt(data.black);
        const whiteBoard = BigInt(data.white);

        // console.log("Num:", blackBoard.toString(2));
        board.updateBoard(blackBoard, whiteBoard);
        
        if (data.black_turn) {
            status.textContent = "Black's turn (●)";
        } else {
            status.textContent = "White's turn (○)";
        }
        count.textContent = `Black-${board.countBlack}   White-${board.countWhite}`

        if (!data.black_turn) {
            pollForBotMove();
        }
    } 
    catch (err) {
        console.error("Move error:", err);
    }
}

let polling = false;

async function pollForBotMove() {
    if (polling) return;
    polling = true;

    const check = async () => {
        const res = await fetch("/state");
        const data = await res.json();

        const blackBoard = BigInt(data.black);
        const whiteBoard = BigInt(data.white);

        board.updateBoard(blackBoard, whiteBoard);
        count.textContent = `Black-${board.countBlack}   White-${board.countWhite}`;

        if (!data.black_turn) {
            setTimeout(check, 500);
        } else {
            blackTurn = data.black_turn;
            status.textContent = blackTurn ? "Black's turn (●)" : "White's turn (○)";
            polling = false;
        }
    };

    check();
}

