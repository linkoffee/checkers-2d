### checkers-2d
`2D checkers, extremely simple version.`

Checkers is a board game for two players where they take turns moving their checkers around a square board and attempting to capture the opponent's checkers. players win by forcing their opponent into a position where they have no more available moves or lose all of their checkers.

![Version](https://img.shields.io/badge/version-0.21-blue)
![License](https://img.shields.io/badge/License-MIT-blue.svg)

### Rules
1. `Playing Field`: Checkers are played on a square 8x8 field, the main game is played on dark squares.
2. `Pieces`: Each player has 12 checkers of one color (usually white and black).
3. `Initial Arrangement`: Checkers are placed on the first three rows on each side, occupying only dark cells.
4. `Game Progress`:
   - Players take turns.
   - Checkers can only move forward one square diagonally.
5. `Eating Checkers`:
   - If the opponent's checker is next to it and there is a free cell diagonally behind it, the player can "eat" the opponent's checker by jumping over it.
   - If several eats are possible in one move, the player must perform them all.
6. `King`:
   - If a player's checker reaches the opponent's last row, it becomes a "King".
   - A King can move both forward and backward, and also jump over the opponent's checkers in any direction.
7. `Win`: The game ends when one of the players cannot make a move (all his checkers are eaten or blocked).
8. `Draw`: If the game is at an impasse and the players cannot make a single move, a draw is declared.

---

#### Stack
- Go 1.22.5
- ebiten 2.7.8

---

#### How to install and play
1. Clone the remote repository and go to it:
```console
git clone https://github.com/linkoffee/checkers-2d.git
```
```console
cd checkers-2d
```
2. Install all dependencies in the main directory:
```console
go mod tidy
```
3. Build the `.exe` file for **Windows** or file without or file without extension for **Linux/macOS**:
```console
go build
```
As a result of executing this command we get an executable file, it will look like this on **Windows**:
```
checkers.go.exe
```
Or like this on **Linux/macOS**:
```
checkers
```
4. Run it on **Windows**:
```console
start checkers.go.exe
```
Or **Linux/macOS**:
```console
./checkers
```
`After that the game file will start, you can rename it as you want, it will not affect its operation.`

<div>
  <img src="https://habrastorage.org/webt/3e/sm/w_/3esmw_pqheskcsvotltc-hnqfmq.png" width="33%" />
  <img src="https://habrastorage.org/webt/9k/gu/2j/9kgu2j982cpb5ev8nxntbfr3uko.png" width="33%" />
</div>

---

Author: [Mikhail Kopochinskiy](https://github.com/linkoffee)
