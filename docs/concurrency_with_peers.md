# Concurrency with Peers, Architecture

we have:

- pieces, files, and lengths. we convert this into a partial piece set, each of which tracks individual blocks and can write to the final file
- we also have a bitfield, of all the pieces we have (which we can broadcast on completion of a piece, with have) - this is less relevant initially, given the above

each peer is a conn, indicating a successful handshake.

for each:

- we exchange bitfields, getting a set of what the peer has

we can then start requesting and receiving pieces. for a single peer this is easy: we send all requests to that peer, throttled a little, and we continuously receive from that peer until all pieces are complete.

for multiple peers, we need to main a goroutine loop for sending and receiving for each peer. these loops need to communicate back to the main thread using channels.

- incoming messages from the center thread tell a peer handler to request a block
- outgoing messages from the peer handler occur when a block is received
- these loops need to check the context to see when done

could there be a single broadcast channel that all peers listen to, or should the center maintain a channel for each peer to send targeted messages through?

peerhandler
    bitfield
    order_chan

pass context and received chan

peer listens for orders and makes requests
might not need an order_chan: the root can just request a piece using a peers conn. the peer handler just has a continuous listener on that conn
