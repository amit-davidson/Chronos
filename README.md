# Chronos

<p align="center">
    <img src="https://i.imgur.com/AhLyxVh.jpeg" width="150" height="225">
</p>

[![made-with-Go](https://github.com/go-critic/go-critic/workflows/Go/badge.svg)](http://golang.org)
[![made-with-Go](https://img.shields.io/badge/Made%20with-Go-1f425f.svg)](http://golang.org)
[![MIT license](https://img.shields.io/badge/License-MIT-blue.svg)](https://lbesson.mit-license.org/)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg?style=flat-square)](http://makeapullrequest.com)
[![amit-davidson](https://circleci.com/gh/amit-davidson/Chronos.svg?style=svg)](https://app.circleci.com/pipelines/github/amit-davidson/Chronos)

Chronos is a static race detector for the Go language written in Go.

## Quick Start:

Download the package

```
go get -v github.com/amit-davidson/Chronos/cmd/chronos
```

Pass the entry point

```
chronos --file <path_to_main> --mod <path_to_module>
```

Help

```
Usage of ./chronos:
  --file string
    	The file containing the entry point of the program
  --mod string
    	Absolute path to the module where the search should be performed. Should end in the format:{VCS}/{organization}/{package}. Packages outside this path are excluded rom the search.
```

## Example:

<p float="left">
    <img src="https://i.imgur.com/LJMP9c2.png" width="245" height="300">
    <img src="https://i.imgur.com/S2mDG7A.png" width="575" height="300">
</p>

## Features:

Support:

- Detects races on pointers passed around the program.
- Analysis of conditional branches, nested functions, interfaces, select, gotos, defers, for loops and recursions.
- Synchronization using mutex and goroutines starts.

Limitations:

- Big programs and external packages. (Due to stack overflow)
- Synchronization using channels, waitgroups, once, cond and atomic.

## Chronos vs go race:

Chronos successfully reports cases where go race fails thanks to his static nature. Mostly because data races appear in unexpected production workloads, which are hard to produce in dev.
In addition, go race is having trouble with short programs where without contrived synchronization the program may exit too quickly.

In contrast, Chronos managed to report only 244/403 = 60.5% of go race test cases. This can be explained by Chronos partial support with Go's features so this number will increase in the future.
Also, it lacked due to his static nature where context/path sensitivity was required.

Therefore, I suggest using both according the strengths and weaknesses of each of the race detectors.

## Credits:

Jan Wen, J., Jhala, R., &amp; Lerner, S. (n.d.). [RELAY: Static Race Detection on Millions of Lines of Code](https://cseweb.ucsd.edu/~lerner/papers/relay.pdf)  
Colin J. Fidge (February 1988). [Timestamps in Message-Passing Systems That Preserve the Partial Ordering"](http://zoo.cs.yale.edu/classes/cs426/2012/lab/bib/fidge88timestamps.pdf)

## More examples:

<p float="left">
    <img src="https://i.imgur.com/NvVWFRf.png" width="230" height="440">
    <img src="https://i.imgur.com/eCNFAX7.png" width=600" height="300">
</p>
<hr style="border:2px solid gray"> </hr>
<p float="left">
    <img src="https://i.imgur.com/app5tBc.png" width="285" height="450">
    <img src="https://i.imgur.com/Lw0LTPp.png" width="545" height="300">
</p>
<hr style="border:2px solid gray"> </hr>