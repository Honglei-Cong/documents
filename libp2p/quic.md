
## QUIC

QUIC (quick udp internet protocol)

1. implemented in user space
2. multiplex-stream
3. tcp fast-open liked, not handshake
4. FEC: 
    1. tcp need timeout to determine if packet lossed, 
    2. QUIC based on FEC, some small packets associated with one checksum


There’s no stream on UDP, this is the big advantage for concurrent transmission and multiplexing stream compared with TCP.
Quic is http over UDP.



