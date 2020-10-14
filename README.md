# Miranda
<p align="center">
    <img src="https://i.imgur.com/AhLyxVh.jpeg" width="150" height="225">
</p>

[![made-with-Go](https://github.com/go-critic/go-critic/workflows/Go/badge.svg)](http://golang.org)
[![made-with-Go](https://img.shields.io/badge/Made%20with-Go-1f425f.svg)](http://golang.org)
[![MIT license](https://img.shields.io/badge/License-MIT-blue.svg)](https://lbesson.mit-license.org/)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg?style=flat-square)](http://makeapullrequest.com)

Miranda is a static race detector for the golang language written in golang.

## Quick Start:
```
go get -u -v github.com/amitdavidson234/Miranda
```

```
Miranda --file <path>
```

## Example:
<p float="left">
    <img style="float: left;" float="left" src="https://i.imgur.com/5td2g2i.png" width="290" height="480">
    <img style="float: left;" float="left" src="https://i.imgur.com/nD9nr9V.png" width="543" height="200">
</p>

## Features:
The project is is still in progress. Therefore, part of the features of the language are supported.

Support:
- Detects races on pointers passed around the program.
- Analysis of conditional branches, nested functions, interfaces and defers.
- Synchronization using mutex and goroutines start.

Limitations:
- Analysis of panics, for loops, goto and recursion.
- Synchronization detection using channels, waitgroups, once, cond.
- Big programs.

## Credits:
Most of my work regarding the detection of the race conditions themselves was based on this paper "RELAY: Static Race Detection on Millions of Lines of Code" by Jan Wen Voung, Ranjit Jhala and Sorin Lerner.