

# Version Vector

The version vector allows the participant to determine if one update preceded another (happended-before), followed it, or if the two updates happened concurrently (and therefore might conflict with each other).

The version vector generates a preorder that tracks the events that precede, and may therefore influence, later updates.

* Initially all vector counters are zero.
* Each time a replica experiences a local update event, it increments its own counter in the vector by one.
* Each two replicas a and b synchronize, they both set the elements in their copy of the vector to the maximum of the element across both counters:  V\_a(x) = V\_b(x) = max(V\_a(x), V\_b(x)).  After synchronization, the two replicas have identical version vectors.


