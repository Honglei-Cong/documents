
## Quic Details

official website: https://www.chromium.org/quic


## Connection Establishment

First time, QUIC client connects to a server

* client send a inchoate client hello (CHLO)
* server sends a rejection (REJ) with the information the client needs to make forward progress

Next time, client send a CHLO

* using the cached credentials from the previous connection to immediately send encrypted requests to the server



## Multiplexing

to solve the HTTP2 (head-of-line-blocking) multiplex issue over TCP

on QUIC, lost packets carrying data for an individual stream generally only impact that specific stream.

Each stream can be immediately dispatched to that stream on arrival, so streams without loss can continue to be reassembled and make forward progress in the application.


## Crypto

* source-address token
    * client’s point of view: opaque byte string
    * server’s point of view: an authenticated-encryption block that contains, at least, the client’s IP address and timestamp by the server
    * receipt of the token by the client is taken as proof of ownership of the IP address in the same way that receipt of a TCP sequence number is.

* Wire protocol
    * full payload of each datagram are authenticated can encrypted once keys have been established.
    * The underlying datagram protocol provides the crypto layer with the means to send reliable, arbitrary sized messages.
    * messages have a uniform, key-value format.
        * key : 32bit tag


