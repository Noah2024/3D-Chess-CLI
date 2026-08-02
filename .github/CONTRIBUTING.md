# Contributing 

This document describes the internal design of **3D-Chess-CLI (3DC)** and the reasoning behind its implementation. 
Before making contributions to the codebase, please read and understand the decisions described here. 
--- 

## Project Goals 

- Implement a complete rules engine for physically three-dimensional chess. 
- Keep the engine modular so individual systems can be tested independently. 
- Keep compiled binary as small reasonable.  
- Optimize for runtime execution speed where possible 
 

---

## Project Structure

```
cmd/
    CLI Entry Points

internal/
    logic called by CLI entrypoints
    game/
        Board creation and loading.
    move/
        Move generation and validation.
        checking/
            Check, checkmate, and pin detection.
        genMovies/
            Move Generation
    ...

util/
    Shared helper packages.

main.go/ 
    Compiled application entry point.
```

---
## Board Representation 

The entire board is represented in memory as a series of length 512 bitmaps. Each bitmap named corresponding to the type of peice its indicies represent. 

![Example Isostructural Screenshot] (images/filestructure.png) 

Because the board is a even 8x8x8 its simple to turn an index to this bitmap (or 'uloc', short for unsigned integer location) into its corresponding position in the x,y, and z directions. Every +1 to a uloc corresponds to a +1 in the x; every +8 corresponds to a +1 to z, and every +64 corresponds to a +1 in the Y.  

Both 3d vectors and uloc's are used throughout the codebase as some tasks are simpler with one as compared to the other.  There are two ways string representation of these bitmaps can be viewed, series of eight 64-bit numbers, or as a length of 512 string of ones and zeros. 

For example, say we have a bitmap representing the position of a pawn, at x:2, y:3, z:4, its corresponding uloc would be (1\*2)+(3\*64)+(4\*8) = 226. Meaning in a 512-length bitmap, the 226th bit would be flipped to a one, and any other pawn of the same team would likewise show up in that same bitmap. 

However, storing many length 512 strings would be very cumbersome for certain pre-generation steps (which are used in move generation). So this same bitmap could be represented by the following string "[0 0 67108864 0 0 0 0 0]". Becuase each layer contains 8x8 or 64 possible spaces each 64 bit 
integer directly corresponds to one layer of the 8x8x8 board.  

Representing the board in this way means doing many calculations are as simple as a few bitwise OR's, AND's and XOR's 

--- 

## Move Pipeline 

Example: 

``` 
Load Board 
↓ 
Determine Bitmap 
↓ 
Generate Raw Moves At Location 
↓ 
Apply Collision Rules 
↓ 
Apply Check Restrictions 
↓ 
Return Legal Moves 

``` 

### Move Generation 

There are two types of move generation, those which take advantage of pre-generated 'dataplanes' and those which have precomputed vectors relative to their starting position. Queens, Bishops, and Rook's all take advantage of data planes. Data planes are pre-generated bitmaps which are dynamically loaded as runtime and correspond to certain planes which cutthrough the 8x8x8 coordinate space of the board.  

For a helpful visual of the dataplanes which correspond to a rook move, check out these 3D views in Desmos3D, [Y](https://www.desmos.com/3d/0grqaurrjl), [Z](https://www.desmos.com/3d/0gdasoxobg), [X](https://www.desmos.com/3d/umhoakbpz4). Using these precomputed data planes, finding all possible moves for a piece is as simple as indexing each of the relevant data planes then computing a bitwise AND between them. (showing why the bitmap representation of board state is so useful).  

For the other types of pieces (King, Knight, and Pawn), while they could've used some fancy way to encode their movements in a bitmap. This was kept simple by precomputing some attack vectors for them, then looping them at runtime. However they're looped using a Range function built into the bitmap package im using, which is SIMD vectorized to keep things as efficient as possible.  


--- 

### Check Detection  


There are two main factors which influence the how this checking system was implemented, 1) The lack of persistent state,  and 2) the necessity to gate a given move if a team is in checkmate. 

Because the runtime environment of go does not support embedding state in the compiled executable itself, caching already computed moves 
becomes very hard, and it becomes vastly simpler to recompute moves every time a move is called.  
Also, because checkmate is dependent on what moves a knig can make, it is necessary to compute all a king's possible moves, every time  
any piece on that team attempts to make a move.  

These two factors combined lead to one, less than an ideal algorithmic situation. Every time we want to determine a checkmate state, we need to run basic movement checks for all pieces currently on the board, both friendly and enemy. This is because we require all possible enemy moves to determine if the king is in checkmate, and we require all possible friendly moves to determine if any piece can save their king from being in check (although friendly move generation is gated behind the wheather the king is in check or not).  

To achieve all features of a standard chess checking system, there are two main functions, kingInDanger, and IsKingInCheck. 

KingInDanger, determines if a king is in check my mathematically verifying if that piece is in line with a king, then stepping in direction away from the king until it finds relevant pieces to determine if it should be pinned. (There are plans to simplify and optimize this later with bitmaps) 

IsKingInCheck, determines InCheck, InCheckmate, ValidKingMoves, SavingKingMoves (those are any moves which could take the king out of check). It does this by first swapping its internal representation of friend and enemy pieces, generating all moves for that team, comparing if the king is in check while doing bitwise operations to restrict king moves. If the king is in check by ONLY one enemy piece, then all friendly moves will be computed to test if any of them intersect with the attack vector that is checking the king. Finally, a checkmate state is determined by seeing if there are any valid king moves or saving king moves.  


As a final note, the outputs are used directly combined in which the previously generated move to compute the moves which are allowed for a user to make. 

--- 

## Piece Logic 

Essentially, just image how a piece moves in chess, and DRAG that up and down, and that's how pieces work. 
I could sit here and try to describe it mathematically but trust me it's far more intuitive than you think.  

## Testing Strategy 

Test everything of major importance.  

As it stands only, the checking and movement generation have test cases associated with them (because it took a very long time to get even those). This is because they are the most important to the proper functioning of the engine. However, more test cases for peripheral and helper functions and packages are planned to be made before the 1.0 release. 

--- 


## Known Limitations 

Currently missing features. 

- Castling 
- En passant 
- Promotion 
- Draw detection 
- AI 

--- 

## Future Improvements 

See readme for roadmap of planned features 

--- 

## Design Philosophy 

General principles used throughout the project. 

- Keep features limited in scope (where not building an OS here) 
- Keep systems independent and modular where possible. 
- Optimize only after correctness. 
- Build clear, robust test cases. 