package main

import (
	// "fmt"
	"log"
	"net"
	"net/rpc"
)

// TODO move common code to shared lib
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

func Dial(host string, rpcname string,
	args interface{}, reply interface{}) bool {
	c, err := rpc.Dial("tcp", host)
	if err != nil {
		return false
	}
	defer c.Close()

	err = c.Call(rpcname, args, reply)
	if err == nil {
		return true
	}

	log.Println(err)
	return false
}

/// end common things that need to move to lib

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
	// Synchronous call
	args := &RegisterArgs{"worker A"}
	var reply RegisterReply
	Dial("127.0.0.1:1234", "Master.Register", args, &reply)
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
