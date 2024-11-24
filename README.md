# Instructions handin 3

1. Start server by running: 
```go run ./Mandatory_Activity_3/Server/server.go```
2. Start client(s) by running ```go run ./Mandatory_Activity_3/Client/client.go``` in separate terminals
3. Clients can send messages by entering a message in the terminal
4. Clients can leave by hitting ```CTRL + C```

# Instructions handin 4

1. Start any number of nodes by running the following command in seperate terminals:
```go run ./Mandatory_Activity_4/Client/client.go```
2. For each node you will be prompted to give an ip for that node, and an ip for the next node in the ring.
3. The program starts running once you create a node with the ip "5052", so this should be the last action performed. This node should be connected to the first node in the ring.

Example of setup with 3 nodes, in the correct order: 
1. Node 1: host ip: "5050", recieving ip: "5051"
2. Node 2: host ip: "5051", recieving ip: "5052"
3. Node 3: host ip: "5052", recieving ip: "5050"


# Source code
The source code for Mandatory Activity 3 is located in ```/Mandatory_Activity_3/```

# Instructions handin 5
1. Start three replica nodes (server nodess) by running the following command in seperate terminals:
   ```go run ./Mandatory_Activity_5/server/server.go```
2. For each replica enter ONE of the following ports: 5050, 5051, 5052
3. Then start any number of clients by running ```go run ./Mandatory_Activity_5/client/Client.go```
4. Then enter an id for each client.
5. Enter 'Bid' to bid on the auction
6. Enter 'Result' to see the status or result of the auction
The auction will end after 30 seconds.

# Source code
The source code for Mandatory Activity 5 is located in ```/Mandatory_Activity_5/```
