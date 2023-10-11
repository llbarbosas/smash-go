<div align="center">

# smash-go

A PoC implementation of a distributed game engine architeture based on SMASH paper

</div>

## About

SMASH (Stackless Microkernel Architecture for SHared environments) is an architecture where a game engine is decomposed in several dynamic and independent software modules interacting via a microkernel-like message bus. The paper can be found [here](https://www.math.unipd.it/~cpalazzi/papers/Palazzi-engine-iscc16.pdf).

## Running

https://github.com/llbarbosas/smash/assets/7810622/862cfaf3-b2c4-4aee-80ec-569dc9ac21f9

A TicTacToe example code is provided in [examples/tictactoe](./examples/tictactoe/). To run it:

```bash
make clean && make modules

# On node 1 terminal
go run cmd/node/node.go

# On node 2 terminal
go run cmd/node/node.go -remotebus :6062 -managment :6063 -link :6060

# On node 1 manager terminal
go run cmd/manager/manager.go
c
s
r

# On node 2 manager terminal
go run cmd/manager/manager.go -addr 127.0.0.1:6063
c
i
r
```
