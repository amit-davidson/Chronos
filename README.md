# Chronos
<p align="center">
    <img src="https://i.imgur.com/AhLyxVh.jpeg" width="150" height="225">
</p>

[![made-with-Go](https://github.com/go-critic/go-critic/workflows/Go/badge.svg)](http://golang.org)
[![made-with-Go](https://img.shields.io/badge/Made%20with-Go-1f425f.svg)](http://golang.org)
[![MIT license](https://img.shields.io/badge/License-MIT-blue.svg)](https://lbesson.mit-license.org/)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg?style=flat-square)](http://makeapullrequest.com)

Chronos is a static race detector for the golang language written in golang.

## Quick Start:
```
go get -u -v github.com/amitdavidson234/Chronos
```

Compile the program and pass it the path to the entry point of the program - main.go
```
chronos --file <main_path>
```

## Example:
<p float="left">
    <img src="https://i.imgur.com/LJMP9c2.png" width="260" height="300">
    <img src="https://i.imgur.com/tWIRIER.png" width="543" height="300">
</p>

## Features:
Support:
- Detects races on pointers passed around the program.
- Analysis of conditional branches, nested functions, interfaces and defers.
- Synchronization using mutex and goroutines start.

Limitations:
- Big programs. (Due to stack overflow)
- Analysis of panics, for loops, goto, recursion and select.
- Synchronization using channels, waitgroups, once, cond and atomic.

## Chronos vs go race:
When compared to go builtin dynamic race detector, Chronos managed to report 244/403 = 60.5% of go race tests. This can be explained by Chronos partial support with the Go's features.
  
In contrast, go race fails to report cases where Chronos succeeds thanks to his static nature. Mostly because race conditions appear in unexpected production workloads which are hard to produce in dev. 


## Credits:
Jan Wen, J., Jhala, R., &amp; Lerner, S. (n.d.). RELAY: Static Race Detection on Millions of Lines of Code. Retrieved from https://cseweb.ucsd.edu/~lerner/papers/relay.pdf

## More examples:
<p float="left">
    <img src="https://i.imgur.com/NvVWFRf.png" width="260" height="450">
    <img src="https://i.imgur.com/VdP7r8B.png" width="543" height="330">
</p>
<hr style="border:2px solid gray"> </hr>
