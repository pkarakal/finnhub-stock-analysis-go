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

### Build & Run 
First, go to [finnhub.io](https://finnhub.io/) and get a free API Key. Save it for later use.

Then to build the application:
1. Fetch the dependencies; this is optional as go will fetch all needed deps during build
    ```shell
    $ go mod download
    ```
2. Compile the project
    ```shell
    $ go build -o finnhub-stock-analysis-go # in the -o flag add whatever name you want.  
    ```
   To cross-compile for Raspberry Pi, just add the environment param `GOARCH=arm` or `GOARCH=arm64` depending on the 
   architecture of your RP. 
   ```shell
   $ GOARCH=arm64 go build -o finnhub-stock-analysis-go
   ```
   If you are compiling on the RP, just run the first command. GO compiler should pick your architecture automatically.
3. Run the binary
    ```shell
    $ ./finnhub-stock-analysis-go  --token <your-finnhub-token> --stocks <stockSymbol> --stocks <stockSymbol>
    ```
