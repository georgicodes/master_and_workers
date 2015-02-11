// Package main provides a sample server program.
package main

import (
	"log"
	"net"
	"net/rpc"
	"os"
	"os/signal"
	"sync"
)

// RegisterArgs ...
type RegisterArgs struct {
	Worker string
}

// RegisterReply ...
type RegisterReply struct {
	OK bool
}

// DoTaskArgs ...
type DoTaskArgs struct {
	Name string
}

// DoTaskReply ...
type DoTaskReply struct {
	OK bool
}

// Dial connects into this server to run some tests.
func Dial(host string, rpcname string, args interface{}, reply interface{}) bool {
	c, err := rpc.Dial("tcp", host)
	if err != nil {
		return false
	}
	defer c.Close()

	if err := c.Call(rpcname, args, reply); err != nil {
		log.Println(err)
		return false
	}

	return true
}

// Master ...
type Master struct {
	l        net.Listener
	wg       sync.WaitGroup
	register chan string
}

// NewMaster creates a new Master value.
func NewMaster() *Master {
	return &Master{
		register: make(chan string),
	}
}

// Register adds a worker to the list of clients that exist.
func (m *Master) Register(args *RegisterArgs, rep *RegisterReply) error {
	log.Println("Registering worker node:", args.Worker)

	// Tell the Master this worker exists.
	m.register <- args.Worker

	rep.OK = true
	return nil
}

// initRPCServer
func (m *Master) initRPCServer() {
	rpc.Register(m)

	var err error
	if m.l, err = net.Listen("tcp", ":1234"); err != nil {
		log.Fatal("listen error:", err)
	}

	m.wg.Add(1)

	// Accept connections on seperate goroutine.
	go func() {
		for {
			conn, err := m.l.Accept()
			if err != nil {
				log.Println("Shutting down Accept handler.")
				// Assume this has an error because we are shutting down.
				m.wg.Done()
				return
			}

			go func() {
				log.Println("Client has connected.")
				rpc.ServeConn(conn)
				log.Println("Client has disconnected.")
				conn.Close()
			}()
		}
	}()
}

// doWork for now takes a registration and uses it as a task.
func (m *Master) doWork() {
	log.Println("Starting to farm out work to workers...")

	m.wg.Add(1)
	go func() {
		for w := range m.register {
			log.Println("Got a worker", w)

			var reply DoTaskReply
			args := DoTaskArgs{
				Name: "task A",
			}

			Dial("127.0.0.1:1235", "Worker.DoTask", &args, &reply)
			log.Println("Result from worker:", reply.OK)
		}

		log.Println("Shutting down doWork handler.")
		m.wg.Done()
	}()
}

// close shutdowns all the goroutines. This does not take an accouting
// of existing client connections at this time.
func (m *Master) close() {
	// Kill the listener.
	m.l.Close()

	// Kill doWork.
	close(m.register)

	// Wait for those to call Done.
	m.wg.Wait()
}

func main() {
	m := NewMaster()
	m.initRPCServer()
	m.doWork()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	<-sigChan

	log.Println("Shutting Down Started")
	m.close()
	log.Println("Shutting Down Completed")
}
