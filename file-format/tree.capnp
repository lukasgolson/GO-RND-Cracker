using Go = import "/go.capnp";
@0x85f1ec173faff1ee;
$Go.package("bk_tree");
$Go.import("olson/bk_tree");


struct Node {
  sequence @0 : List(Int8);
  # The sequence associated with this node, representing a potential signature.
  seed @1 : Int32;
  # The seed value used to generate the pseudorandom sequence.
  children @2 : List(Edge);
  # List of edges connecting to child nodes representing similar signatures.
}

struct Edge {
  distance @0 : Int8;
  # Edit distance between parent node's sequence and child node's sequence (signature).
  node @1 : Node;
  # Child node connected by this edge, representing a similar signature.
}
