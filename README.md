# FinnHub stocks analysis  
This project is implemented as a part of a homework exercise for [077] - Real Time Embedded Systems
course of ECE Department, AUTh.

This fetches stock information changes as they happen in real time, using websockets, runs some analysis
and persists the information as it becomes available. It should consume the least cpu it can as it should
be able to be deployed in a microcontroller. 


## Getting Started

### Prerequisites
1. Go
To install them on variant Linux distributions follow the instructions below

#### Fedora
```shell
$ sudo dnf upgrade --refresh # updates installed packages and repositories metadata
$ sudo dnf install golang # or follow official instructions https://go.dev/doc/install
```

#### Ubuntu
```shell
$ sudo apt-get update && sudo apt-get upgrade # updates installed packages and repositories metadata
$ sudo apt-get install golang # or follow official instructions https://go.dev/doc/install
```


[More info incoming as the project progresses]
