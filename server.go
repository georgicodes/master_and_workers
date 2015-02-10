package main

import (
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

type Master struct {
	l               net.Listener
	isAlive         bool
	doneChan        chan bool
	registerChannel chan string
}

// Register workers
func (m *Master) Register(args *RegisterArgs, rep *RegisterReply) error {
	log.Printf("Registering worker node: %s\n", args.Worker)
	go func() { // TODO: do i need a goroutine here, clients don't seem to get a response w/o it.
		m.registerChannel <- args.Worker
	}()
	rep.OK = true
	return nil
}

func InitMaster() *Master {
	m := new(Master)
	m.isAlive = true
	m.doneChan = make(chan bool)
	m.registerChannel = make(chan string)
	return m
}

func (m *Master) initRPCServer() {
	rpc.Register(m)

	listener, e := net.Listen("tcp", ":1234")
	if e != nil {
		log.Fatal("listen error:", e)
	}
	m.l = listener

	// accept connections on seperate thread.
	go func() {
		for m.isAlive {
			conn, err := m.l.Accept()
			if err == nil {
				go func() {
					log.Println("serving request")
					rpc.ServeConn(conn)
					conn.Close()
				}()
			} else {
				log.Println("errors in go routine")
			}
		}
	}()
}

// Ranges over lists of workers and asks each of them to DoWork
func (m *Master) doWork() {
	log.Println("Starting to farm out work to workers...")

	for w := range m.registerChannel {
		log.Printf("Got a worker %s", w)
		// Synchronous call
		args := &DoTaskArgs{"task A"}
		var reply DoTaskReply
		Dial("127.0.0.1:1235", "Worker.DoTask", args, &reply)
		log.Printf("Result from worker: %v", reply.OK)
	}

	// program will exit when DoneChan = true
	<-m.doneChan
}

func main() {
	m := InitMaster()
	m.initRPCServer()
	m.doWork()
}
