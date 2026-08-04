# 3D-Chess-CLI

![License Badge](https://img.shields.io/badge/GPL%203.0%20-%20red?logo=gplv3&label=license&link=https%3A%2F%2Fwww.gnu.org%2Flicenses%2Fgpl-3.0.en.html)
![Go Version Badge](https://img.shields.io/badge/v1.25.5-blue?logo=go&label=go)
![GitHub Release Badge](https://img.shields.io/github/v/release/Noah2024/3D-Chess-CLI?include_prereleases&link=https%3A%2F%2Fgithub.com%2FNoah2024%2F3D-Chess-CLI%2Freleases)
![GitHub Repo Size Badge](https://img.shields.io/github/repo-size/Noah2024/3D-Chess-CLI?color=darkyellow)


3D-Chess-CLI (3DC) is just as the name implies a Command-Line-Interface built built to make, change, and validate 3 Dimensions of play for a game of chess. 

## What is meant by 3 Dimensions?

3DC is not meant for a graphically three dimensional game, but A LITTERLY three dimensional game of chess. 
The standard board size is 8x8x8 (just a regular board stacked ontop of itself 8 times).
Pieces can move, attack, check, and deliver checkmate in all three dimensions. 

## Current State

As it stands much of the basics of chess work and work well, HOWEVER certian features are not implmented and some won't be at all.
For starters there will be no castling. In a normal game of chess castling is very often a very good way to protect your king, in 3DC however its not.
Due to the phyical laws of nature, in 3DC your king could still be attacked 14 differnet ways even after castling, so I feel its really not worth it. 

Additionally pawns CURRENTlY cannont promot, move twice from their starting square, or make an en-pessent move (for the 5 people who actually en-pessent).
This is becuase the pawn is a stupid, ridicious piece, whos very existence is predicated on being differnet in every concivable way from everyone else.
However I do have plans to implment these before a 1.0 release. 

## Building From Source

```
make build
```

## Usage

Use '-h' at any time and on any command for extra help

```
board/ - commnads relating to board
    view - to view a certian board (accepts optional paramter of layers UPPERCASE letter)

game/ - commands related to managing games
    delete - deletes a given name by either name or listed number
    list - lists all games in saved game directory with both name and number
    load - loads a given game from saved game direcotry using either name or number
    new - creates new game in default start state, destroys previous game
    save - saves CurrentGame to saved game directory with given name

move/ 
    move - takes two arguments, locationFrom and locationToo, expected in format a1A, and h1H (x, y, z)

debug/ - commands related to debug (if using please read implemnetation)
   bmstring
   dataPlanes
   moves
   uintvec
   vecuint
```

*Tip* The standard game starts on layer C, so to ground yourself try doing ```3DC board view C``` after creating a new game

## Road Map

### Final Checklist before 1.0
- Standardize debug logs & improve storage
- Clean Up/Refactor and modularize Checking logic
- Change view for terminal 
- Actually Integrate metadata
- Ensure UintToVec and VecToUint are really inverses
- Systems Test Using the full move command
- Turn Tracking
- Upper Level API Integration
 - And test cases for said API 
- pawn double move
- pawn promotion
- en pessant 🤮
- Time trials/stress test
- Optimize w/ concurrent goroutine

### Further Future 
- Compilation to Web Assembly
- Web app to play games
- Docker Image for self hosting move validation engine
