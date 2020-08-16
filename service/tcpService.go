package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/Riften/libp2p-playground/util"
	"io"
	"log"
	"math/rand"
	"net"
	"sync"
	"time"
)

const DefaultTCPPort = 10000

type TCPService struct {
	//port         int
	ctx 		 context.Context
	senders		map[string] *TcpSender
	receivers	map[string] *TcpReceiver

	lock sync.Mutex
}

type TcpSender struct {
	Ip string
	Recorder *util.SpeedRecorder
	canceler context.CancelFunc
	port int

	active bool
}

type TcpReceiver struct {
	Remote   string
	Recorder *util.SpeedRecorder
	conn     net.Conn

	active bool
}

func (t *TCPService) GetSender(ip string) (*TcpSender, bool) {
	t.lock.Lock()
	defer t.lock.Unlock()
	tmpSender, ok := t.senders[ip]
	return tmpSender, ok
}

func (t *TCPService) GetReceiver(addr string) (*TcpReceiver, bool) {
	t.lock.Lock()
	defer t.lock.Unlock()
	tmpSender, ok := t.receivers[addr]
	return tmpSender, ok
}

// Return all the ip sending to.
// The non-active senders would be removed.
func (t *TCPService) SendList() ([]string) {
	t.lock.Lock()
	defer t.lock.Unlock()
	resList := make([]string, 0)
	removeList := make([]string, 0)
	for k, sender := range t.senders {
		if sender.active {
			resList = append(resList, k)
		} else {
			log.Println("sender for ", k, " not active, removed")
			removeList = append(removeList, k)
		}
	}

	for _, k := range removeList {
		delete(t.senders, k)
	}
	return resList
}

// Return all the address receive from.
// The non-active receivers would be removed.
func (t *TCPService) ReceiveList() ([]string) {
	t.lock.Lock()
	defer t.lock.Unlock()
	resList := make([]string, 0)
	removeList := make([]string, 0)

	for k, receiver := range t.receivers {
		if receiver.active {
			resList = append(resList, k)
		} else {
			log.Println("receiver from ", k, " not active, removed")
			removeList = append(removeList, k)
		}
	}

	for _, k := range removeList {
		delete(t.receivers, k)
	}
	return resList
}

func NewTCPService(ctx context.Context) *TCPService {
	return &TCPService{
		//port:         port,
		ctx:		  ctx,
		senders:	  make(map[string]*TcpSender),
		receivers: 	  make(map[string]*TcpReceiver),
	}
}

// Run this in seperate routine
func (t *TCPService) StartListen(port int) error {
	listen_sock, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", port))
	if err != nil {
		log.Println("Error when listen tcp port ", port, ": ", err)
		return err
	}

	go func() {
		for {
			new_conn, err := listen_sock.Accept()
			if err != nil {
				log.Println("Error when accept new connection: ", err)
				continue
			}
			log.Println("New connection from ", new_conn.RemoteAddr().String())
			go t.handleRecv(new_conn)
		}
	}()
	return nil
}

func (t *TCPService) StartSend(ip string, port int) error {
	t.lock.Lock()
	defer t.lock.Unlock()
	_, ok := t.senders[ip]
	if ok {
		log.Println("sender already exists.")
		return errors.New("sender for "+ip+" already exists")
	}
	newCtx, canceler := context.WithCancel(t.ctx)
	newSender := &TcpSender{
		Ip:       ip,
		Recorder: util.NewSpeedRecorder(newCtx, 10),
		canceler: canceler,
		port:     port,
		active:   false,
	}

	t.senders[ip] = newSender
	log.Println("sender add: ", ip)
	go newSender.Recorder.Start()
	go newSender.startSend()
	return nil
}

func (t *TCPService) StopSend(ip string) error {
	t.lock.Lock()
	defer t.lock.Unlock()
	sender, ok := t.senders[ip]
	if !ok {
		return errors.New("sender for " + " not exists")
	}
	sender.canceler()
	delete(t.senders, ip)
	return nil
}

// ==================== Private
func (t *TCPService)handleRecv(conn net.Conn) {
	t.lock.Lock()
	defer t.lock.Unlock()
	remoteAddr := conn.RemoteAddr().String()
	_, ok := t.receivers[remoteAddr]
	if ok {
		log.Println("Error when handle receive: receiver already exists for ", remoteAddr)
		err := conn.Close()
		if err != nil {
			log.Println("Error when close receive connection: ", err)
		}
		return
	}
	newReceiver := &TcpReceiver{
		Remote:   remoteAddr,
		Recorder: util.NewSpeedRecorder(t.ctx, 10),
		conn:     conn,
		active:   false,
	}
	t.receivers[remoteAddr] = newReceiver
	log.Println("receiver add: ", remoteAddr)
	go newReceiver.Recorder.Start()
	go newReceiver.startReceiver()
}

func (s *TcpSender) startSend()  {
	rander := rand.New(rand.NewSource(time.Now().UnixNano()))
	batchSize := 1024*1024	// 1024kb, 1MB
	batch := make([]byte, batchSize)
	var writeSize int
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", s.Ip, s.port))
	if err != nil {
		log.Println("Error when dial tcp: ", err)
		return
	}
	s.active = true

	defer func() {
		conn.Close()
		s.active = false
	}()
	for {
		_, err = rander.Read(batch)
		if err != nil {
			log.Println("Error when generate random bytes: ", err)
			return
		}
		writeSize, err = conn.Write(batch)
		if err != nil {
			log.Println("Error when write to connection: ", err)
			return
		}
		s.Recorder.AddStamp(time.Now(), writeSize)
	}
}

func (r *TcpReceiver) startReceiver() {
	r.active = true
	defer func() {
		r.active = false
	}()
	bufSize := 1024 * 1024 //1MB
	readBuf := make([]byte, bufSize)
	var readSize int
	var err error
	for {
		readSize, err = io.ReadFull(r.conn, readBuf)
		if err != nil {
			log.Println("Error when read from connection: ", err)
			return
		}
		r.Recorder.AddStamp(time.Now(), readSize)
	}
}
