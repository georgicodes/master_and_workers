// Don't know enough about all this code yet. Just cleaned up some idiomatic issues.
// Made some changes but did not try to build them.

// Always code negative path inside if statements. Keep positive path outside.
// Not sure I like these single field structs. I am assuming these will grow.
// Add your proper comments as you code.
// Make sure you are also go vet and golint on every save.

package main

import (
	// "fmt"
	"log"
	"net"
	"net/rpc"
)

// if and only if these types work together.
type (
	RegisterArgs struct {
		Worker string
	}

	RegisterReply struct {
		OK bool
	}

	DoTaskArgs struct {
		Name string
	}

	DoTaskReply struct {
		OK bool
	}
)

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

	return false
}

type Worker struct {
	l        net.Listener
	isAlive  bool
	doneChan chan bool
}

func (w *Worker) DoTask(args *DoTaskArgs, rep *DoTaskReply) error {
	log.Println("Asked by master to do task:", args.Name)
	rep.OK = true
	return nil
}

func (w *Worker) registerWithMaster() {
	var reply RegisterReply
	args := RegisterArgs{
		Worker: "worker A",
	}
	
	Dial("127.0.0.1:1234", "Master.Register", &args, &reply)
	log.Println("Result from master:", reply.OK)
}

func (w *Worker) registerAndWaitForTasks() {
	rpc.Register(w)

	var err error
	w.l, err = net.Listen("tcp", ":1235")
	if err != nil {
		log.Fatal("listen error:", e)  // Why not handle an error?
	}
	w.registerWithMaster()

	// accept connections on seperate thread.
	go func() {
		for w.isAlive {
			conn, err := w.l.Accept()
			if err != nil {
				log.Println("errors in go routine")
				continue
			}
			
			go func() {
				log.Println("serving requests from master")
				rpc.ServeConn(conn)
				conn.Close()
			}()
		}
	}()
}

func InitWorker() *Worker {
	return &Worker {
		isAlive: true,
		doneChan:  make(chan bool),
	}
}

func main() {
	w := InitWorker()
	w.registerAndWaitForTasks()

	// program will exit when DoneChan = true
	<-w.doneChan
}
