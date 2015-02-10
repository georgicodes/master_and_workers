package main

import (
	"log"
	"net"
	"net/rpc"
)

// Same as before, only if these work together maybe group them.

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

func Dial(host string, rpcname string, args interface{}, reply interface{}) bool {
	c, err := rpc.Dial("tcp", host)
	if err != nil {
		return false
	}
	defer c.Close()

	if err := c.Call(rpcname, args, reply) err != nil {
		log.Println(err)
		return false
	}

	return true
}

type Master struct {
	l               net.Listener
	isAlive         bool
	doneChan        chan bool
	registerChannel chan string
}

// Register workers
func (m *Master) Register(args *RegisterArgs, rep *RegisterReply) error {
	log.Println("Registering worker node:", args.Worker)
	go func() { // TODO: do i need a goroutine here, clients don't seem to get a response w/o it.
		m.registerChannel <- args.Worker
	}()
	rep.OK = true
	return nil
}

func InitMaster() *Master {
	return &Master{
		isAlive: true,
		doneChan: make(chan bool),
		registerChannel: make(chan string),
	}
}

func (m *Master) initRPCServer() {
	rpc.Register(m)

	var err error
	m.l, err = net.Listen("tcp", ":1234")
	if e != nil {
		log.Fatal("listen error:", e)
	}

	// accept connections on seperate thread.
	go func() {
		for m.isAlive {
			conn, err := m.l.Accept()
			if err != nil {
				log.Println("errors in go routine")
				continue
			}
			
			go func() {
				log.Println("serving request")
				rpc.ServeConn(conn)
				conn.Close()
			}()
		}
	}()
}

func (m *Master) doWork() {
	log.Println("Starting to farm out work to workers...")

	for w := range m.registerChannel {
		log.Println("Got a worker", w)
		
		var reply DoTaskReply
		args := DoTaskArgs{
			name: "task A"
		}
		
		Dial("127.0.0.1:1235", "Worker.DoTask", &args, &reply)
		log.Println("Result from worker:", reply.OK)
	}

	// program will exit when DoneChan = true
	<-m.doneChan
}

func main() {
	m := InitMaster()
	m.initRPCServer()
	m.doWork()
}
