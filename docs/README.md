# PSI-STATS Documentation

[PSI-STATS](https://eprint.iacr.org/2020/623.pdf) is a Private Set Intersection statistics protocol that allows
two parties to compute sample statistics such as the mean on the intersection of their respective sets.

PSI-STATS is a secure multi-party protocol that allows to compute the intersection of two sets and statistics
on said sets.

Unlike protocols based on static hash functions or Diffie-Hellman, PSI-STATS uses a mixture of hash functions
and modular exponentiation to generate successive commitments to identifier values (items in the set)

The computations are done using homomorphic encryption schemes such as [Paillier](https://en.wikipedia.org/wiki/Paillier_cryptosystem)
to compute summary statistics over the intersection.
