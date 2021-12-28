# Protocol Definitons

This document holds definitions of all protocols supported by `psistats`.

## Mean-Intersection

The Mean-Intersection protocol computes the mean of the set intersection between two parties A and B.
Consider Alice and Bob are two parties with a set of key values (k,v) describing user statistics *t_i*
over a set of users **U_i**.

Consider Alice's set S_a (a[i],t[i]) and S_b(b[i],t[i]), since the intersection of the identifiers
a[i] INTERSECT b[i] can be done in the SMPC (A and B learn nothing more than the intersection)
it is feasible for both A and B to compute the mean of the intersection using a homomorphic scheme,
in our case we use the Paillier scheme.
