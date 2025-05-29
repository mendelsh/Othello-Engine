class Board {
    constructor(containerId) {
      this.container = document.getElementById(containerId);
      this.cells = [];
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
        for (let i = 0; i < 64; i++) {
            const blackBit = (blackBoard >> BigInt(i)) & 1n;
            const whiteBit = (whiteBoard >> BigInt(i)) & 1n;
        
            // Clear previous disc if any
            this.cells[i].innerHTML = '';
        
            if (blackBit === 1n) {
                const disc = document.createElement('div');
                disc.classList.add('disc', 'black');
                this.cells[i].appendChild(disc);
            } else if (whiteBit === 1n) {
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

        board.updateBoard(blackBoard, whiteBoard);
        
        if (data.black_turn) {
            status.textContent = "Black's turn (●)";
        } else {
            status.textContent = "White's turn (○)";
        }
    } 
    catch (err) {
        console.error("Move error:", err);
    }
}


