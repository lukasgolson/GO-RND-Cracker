# Algorithm: Insertion into a BK-tree

## Input:
- `t`: The BK-tree.
- `d(uv)`: Weight assigned to an arc (u, v).
- `w(u)`: Word assigned to a node u.
- `d`: The discrete metric used by 't' (e.g., Levenshtein distance).
- `w`: The element to be inserted into 't'.

## Output:
- The node in 't' corresponding to 'w'.

## Procedure:
1. If 't' is empty:
   1.1 Create a root node 'r' in 't'.
   1.2 Set 'w(r)' to 'w'.
   1.3 Return 'r'.

2. Set 'u' as the root of 't'.

3. While 'u' exists:
   3.1 Calculate 'k' as 'd(w(u), w)'.
   3.2 If 'k' is equal to 0:
   3.2.1 Return 'u'.

   3.3 Find 'v,' the child of 'u' such that 'd(uv)' equals 'k.'
   3.4 If 'v' is not found:
   3.4.1 Create the node 'v.'
   3.4.2 Set 'w(v)' to 'w.'
   3.4.3 Create the arc '(u, v)' with weight 'd(uv)' set to 'k.'
   3.4.4 Return 'v.'

   3.5 Set 'u' to 'v'.
