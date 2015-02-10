package main

import (
	// "fmt"
	"log"
	"net"
	"net/rpc"
)

// TODO move to shared lib
type RegisterArgs struct {
	Worker string
}

type RegisterReply struct {
	OK bool
}

type DoTaskArgs struct {
	Name string
}

type DoTaskReply struct {
	OK bool
}

type Worker struct {
	l        net.Listener
	isAlive  bool
	doneChan chan bool
}

func (w *Worker) DoTask(args *DoTaskArgs, rep *DoTaskReply) error {
	log.Printf("Asked by master to do task: %s", args.Name)
	rep.OK = true
	return nil
}

func (w *Worker) registerWithMaster() {
	c, err := rpc.Dial("tcp", "127.0.0.1:1234") // TODO: no magic strings for master settings
	if err != nil {
		log.Fatal("dialing:", err)
	}

	// Synchronous call
	args := &RegisterArgs{"worker A"}
	var reply RegisterReply
	err = c.Call("Master.Register", args, &reply)
	defer c.Close()
	if err != nil {
		log.Fatal("Error connecting to remote: %s", err)
	}
	log.Printf("Result from master: %v", reply.OK)
}

func (w *Worker) registerAndWaitForTasks() {
	rpc.Register(w)

	listener, e := net.Listen("tcp", ":1235")
	if e != nil {
		log.Fatal("listen error:", e)
	}
	w.l = listener
	w.registerWithMaster()

	// accept connections on seperate thread.
	go func() {
		for w.isAlive {
			conn, err := w.l.Accept()
			if err == nil {
				go func() {
					log.Println("serving requests from master")
					rpc.ServeConn(conn)
					conn.Close()
				}()
			} else {
				log.Println("errors in go routine")
			}
		}
	}()
}

func InitWorker() *Worker {
	w := new(Worker)
	w.isAlive = true
	w.doneChan = make(chan bool)
	return w
}

func main() {
	w := InitWorker()
	w.registerAndWaitForTasks()

	// program will exit when DoneChan = true
	<-w.doneChan
}
