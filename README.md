# 3D-Chess-CLI

![Static Badge](https://img.shields.io/badge/GPL%203.0%20-%20red?logo=gplv3&label=license&link=https%3A%2F%2Fwww.gnu.org%2Flicenses%2Fgpl-3.0.en.html))
![Static Badge](https://img.shields.io/badge/v1.25.5-blue?logo=go&label=go)

Simple idea, I want to make 3d chess, I've made an old shitty version in python, but that was very specific needed big libraries and took alot of time to run.

SO Im now using this project to a) work on making that dream a reality 2) learn how to structure a go project, and III) have a bit of fun along the way. Once this shit is actually in a workable state, i'll get on some actual documentation and work on a v1 release

## Final Checklist before 0.5
- touch up code
- documentation
- contribution guide 


## Final Checklist before 1.0
- Clean Up/Refactor and modularize Checking logic
- Actually Integrate metadata
- Ensure UintToVec and VecToUint are really inverses
- Systems Test Using the full move command
- Turn Tracking
- Upper Level API Integration
 - And test cases for said API 
- pawn double move
- pawn promotion
- Time trials/stress test
- Optimize w/ concurrent goroutine
