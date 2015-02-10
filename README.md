Master + Workers in Go
======================

This is a simple master and worker app in Go. Communication is done via RPC. The workers each register themselves with the master and then wait to be given tasks to do.

```bash
# Run the server
go build server.go
./server

# run the clients
go build client.go
./client (multiple can be run)
```