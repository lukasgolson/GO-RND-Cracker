# BK-Tree Lookup Algorithm

## Overview

The BK-Tree Lookup Algorithm is used to find the closest element to a searched element 'w' within a BK-tree. This algorithm efficiently explores the tree, taking advantage of the BK-tree's structure and the triangle inequality as a cut-off criterion to improve search performance.

## Input

- **t:** The BK-tree, representing a set of elements.
- **d:** The corresponding discrete metric, e.g., the Levenshtein distance, used for measuring element similarity.
- **w:** The element being searched for.
- **d_max:** The maximum allowed distance between the best match and 'w,' defaulting to positive infinity (∞).

## Output

- **w_best:** The closest element to 'w' stored in 't' according to 'd,' or '⊥' (perpendicular symbol) if not found.

## Algorithm

1. If the BK-tree 't' is empty:
    - Return '⊥' to indicate no match was found.

2. Initialize a set 'S' to hold nodes for processing and insert the root of 't' into 'S.'
    - Initialize (w_best, d_best) as (⊥, d_max).

3. While 'S' is not empty:
   a. Pop an arbitrary node 'u' from 'S.'
    - Calculate d_u = d(w, w_u), where w_u is the element represented by node 'u.'
      b. If d_u is less than d_best:
    - Update (w_best, d_best) to (w_u, d_u).
      c. For each egress-arc (u, v):
    - If |d_uv - d_u| is less than d_best (cut-off criterion), insert 'v' into 'S.'

4. Return 'w_best' as the closest element found in 't' based on the specified metric 'd.'
