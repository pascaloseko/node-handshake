# BitCoin Node Handshake

This is an implementation of a tcp network handshake with a publicly available bitcoin node

## Choice of Bitcoin Node to use
I decided to go for [bitcoin core](https://github.com/bitcoin/bitcoin) as the primary choice because

1. It is widely used and trusted since it's being maintained by a large community of developers it has undergone extensive tests and review.
2. It provides scalability and reliablity due to regular updates to address bugs, security vulnerabilities and performance improvements.
3. It has an active development, ensuring that it stays up-to-date with the latest developments in the Bitcoin ecosystem.
4. It has a comprehensive documentation, making it easier to understand and work with for tasks such as network communication.
5. It is fully compatible with the Bitcoin network, ensuring that the handshake process will be consistent with other Bitcoin nodes.

## Design
The Go code serves as the client, and the Bitcoin node implementation serves as the server during the network handshake process.

## Requirements
- Both the target node and the handshake code should compile at least on Linux.
- The solution has to perform a full protocol-level (post-TCP) handshake with the target node.
- The solution can not use the node implementation as a dependency.
- The submitted code can not use existing P2P libraries for the handshake.

## Implementation

1. Download and run the bitcoin node server
    - Follow the instructions provided in the Bitcoin Core documentation to download and install the software on the Linux system. Ensure that it compiles and runs smoothly on Linux.
2. Configure Bitcoin Core
    - Modify the Bitcoin Core configuration file (bitcoin.conf) to enable listening for incoming connections (listen=1). Ensure that the node is reachable from the network by configuring any necessary firewall or router settings.
3. Run Bitcoin Node Server.
    - Start the Bitcoin Core server on the Linux machine using the bitcoind command. Monitor the server's output to ensure it starts without errors and begins synchronizing with the Bitcoin network.
4. Implement Client Handshake:
    - Develop the Go code to perform the handshake with the Bitcoin node. Ensure that the code compiles and runs on Linux without dependencies on the Bitcoin Core implementation.
6. Test Handshake Connection
    - Once the unit tests pass, conduct integration testing by establishing a handshake connection with the running Bitcoin Core server. Verify that the handshake completes successfully and that the client can communicate with the node.

## How to run it locally
- Run Bitcoin Node Server
    ```
    ./bitcoind -regtest -daemon
    ```
- Build the Go client
    ```
    go build -o handshake
    ```
- Run the go client
    ```
    ./handshake
    ```

## Problems [UPDATE this is fixed]
- I was not able to correctly send the handshake payload to the server hence the client terminated when running it 
