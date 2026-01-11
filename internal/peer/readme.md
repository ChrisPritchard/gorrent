# Peer communication

This is a state machine.

peers start off choked

1. send interested, peer responds with unchoked
2. send bitfield, receive bitfield, indicating what they have and what you have
3. send request messages
4. parse responses, receiving piece messages with data
5. once a piece is completed, and the hash is checked, a 'have' message is sent
6. if the peer responds with choke then stop sending requests
7. if a request has been sent to a peer, but it was received from somewhere else, a cancel message can be sent to that peer

with multiple peers, this sort of communication is all managed at once

choked, interested, bitfield per peer

for each peer we need to listen to responses

process:

- maintain up to five requests at a time
- we issue more requests as downloads are received
- so we need a shared state, with the count of active requests, and which can be updated by per-peer managers
- each peer has its own state machine
- this checks state, and sends unchoke, interested messages or parts if requested
- it receives messages and updates the global shared state as appropriate
