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
