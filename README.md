# psistats

## Project structure

- `cmd/` : Implements a lightweight sandbox for executing the protocol.
- `data/` : Synthetic datasets for benchmarks and tests.
- `docs/` : Project and Protocol docs.
- `internal/` : Implementation internals.
- `pkg/` : Implementation packages with an external API.
    - `crypto/`: Provides an easy to mis-use interface to the cryptographic building blocks.
    - `net/` : Provides a 2-Party communication channel using Websockets and HTTP/2.
    - `proto/` : Provides an API to protocol definition and constants.
- `scripts/` :  Useful scripts for generating synthetic datasets for testing and simulations.
- `tests/` : Integration tests.

## Paper details

We define protocol parameters in pseudocode and describe core algorithms and building blocks.

```sh

a[i],b[i] # identifiers
t[i]  # statistic associated with identifier i.
A.X = [a[i]] # party A holds set X
B.Y = [b[i],t[i]] # party B holds set Y
Y_ = [b[i]]
E = (Enc,Dec) # agreed upon instance of Paillier
h = (Sha256) # sha256 hash function
G # group of large order

```