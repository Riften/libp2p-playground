package service

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	"io"
	"math/rand"
	"time"

	//"github.com/libp2p/go-libp2p-core/protocol"
	"github.com/orcaman/concurrent-map"
	"log"
)

const SpeedTestProtocol = "speedtest/1.0"

type streamSession struct {
	stream network.Stream
	remotePeer peer.ID
	rw *bufio.ReadWriter
	sender *sendManager
	recver *recvManager
}

type SpeedService struct {
	tasks cmap.ConcurrentMap
	host host.Host
	printer *speedPrinter
	ctx context.Context
}

func NewSpeedService(h host.Host, ctx context.Context,) *SpeedService {
	newService := &SpeedService{
		tasks:   cmap.New(),
		host:    h,
		printer: &speedPrinter{
			sendSpeedChan: make(chan *record, 5),
			recvSpeedChan: make(chan *record, 5),
		},
		ctx:     ctx,
	}
	return newService
}

func (s *SpeedService) Start() {
	go s.printer.run()
	s.host.SetStreamHandler(SpeedTestProtocol, s.handleStream)
}

func (s *SpeedService) handleStream(stream network.Stream) {
	fmt.Println("SpeedTest: New stream got.")
	rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))
	sess := &streamSession{
		stream:     stream,
		remotePeer: stream.Conn().RemotePeer(),
		rw:         rw,
	}
	s.tasks.Set(stream.Conn().RemotePeer().Pretty(), sess)
	newCtx, canceler := context.WithCancel(s.ctx)
	sess.recver = &recvManager{
		ctx:    newCtx,
		cancel: canceler,
		rw:     rw,
		out:    s.printer.recvSpeedChan,
	}
	go sess.recver.run()
}

func (s *SpeedService) StartSend(pid string) error {
	var sess *streamSession
	tmp, ok := s.tasks.Get(pid)
	if ok {
		sess = tmp.(*streamSession)
		log.Println("Session for ", pid, " already exists.")
	} else {
		peerId, err := peer.Decode(pid)
		if err != nil {
			log.Println("Error when decode peer id: ", err)
			return err
		}
		stream, err := s.host.NewStream(s.ctx, peerId, SpeedTestProtocol)
		if err != nil {
			log.Println("Error when create stream with "+pid+": ", err)
			return err
		}
		rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))
		sess = &streamSession{
			stream:     stream,
			remotePeer: peerId,
			rw:         rw,
			//sender:     ,
		}
		s.tasks.Set(pid, sess)
	}
	if sess.sender != nil {
		log.Println("Sender already exists")
		return errors.New("sender already exists")
	}
	newCtx, cancel := context.WithCancel(s.ctx)
	sender := &sendManager{
		ctx:    newCtx,
		cancel: cancel,
		rw:     sess.rw,
		out:    s.printer.sendSpeedChan,
	}
	sess.sender = sender
	go sender.run()
	return nil
}

type record struct {
	//Peer string `json:"peer"`
	Speed int64 `json:"speed"`
	AvgSpeed int64 `json:"avgSpeed"`
}

// - Cancel the send routine
// - Keep the send speed
// - Through out the send speed
type sendManager struct {
	ctx context.Context
	cancel context.CancelFunc
	rw *bufio.ReadWriter
	out chan *record
}

type recvManager struct {
	ctx context.Context
	cancel context.CancelFunc
	rw *bufio.ReadWriter
	out chan *record
}

func (m *sendManager) run() {
	timeStamp := time.Now()
	var newTimeStamp time.Time
	var duration int64 //in millsecond
	rander := rand.New(rand.NewSource(timeStamp.UnixNano()))
	batchSize := 1024*1024	// 1024kb, 1MB
	batch := make([]byte, batchSize)
	var totalMs int64 = 0
	var totalByte int64 = 0

	for {
		select {
		case <-m.ctx.Done():
			log.Println("send end: ", m.ctx.Err())
			return
		default:
			_, err := rander.Read(batch)
			if err != nil {
				log.Println("Error when generate random bytes: ", err)
				return
			}
			writeBytes, err := m.rw.Write(batch)
			writeBytes64 := int64(writeBytes)
			if err != nil {
				log.Println("Error when send bytes: ", err)
				return
			}
			newTimeStamp = time.Now()
			duration = newTimeStamp.Sub(timeStamp).Milliseconds()
			totalMs += duration
			totalByte += writeBytes64
			if duration == 0 {
				continue
			}
			m.out <- &record{
				//Peer:     "",
				Speed:    (1000*writeBytes64/1024)/duration,
				AvgSpeed: (1000*totalByte/1024)/totalMs,
			}
			timeStamp = newTimeStamp
		}
	}
}

func (m *recvManager) run() {
	timeStamp := time.Now()
	var newTimeStamp time.Time
	var duration int64 //in millsecond

	batchSize := 1024*1024	// 1024kb, 1MB
	batch := make([]byte, batchSize)


	var totalMs int64 = 0
	var totalByte int64 = 0

	for {
		select {
		case <-m.ctx.Done():
			log.Println("send end: ", m.ctx.Err())
			return
		default:
			readBytes, err := io.ReadFull(m.rw, batch)
			 //:= m.rw.Read(batch)
			readBytes64 := int64(readBytes)
			if err != nil {
				log.Println("Error when read from stream: ", err)
				return
			}

			newTimeStamp = time.Now()
			duration = newTimeStamp.Sub(timeStamp).Milliseconds()
			totalMs += duration
			totalByte += readBytes64
			if duration == 0 {
				continue
			}
			m.out <- &record{
				//Peer:     "",
				Speed:    (1000*readBytes64/1024)/duration,
				AvgSpeed: (1000*totalByte/1024)/totalMs,
			}
			timeStamp = newTimeStamp
		}
	}
}

type speedPrinter struct {
	sendSpeedChan chan *record
	recvSpeedChan chan *record
}
func printSpeed(up int64, avgUp int64, down int64, avgDown int64) {
	fmt.Printf("上行 %6dkb/s 平均 %6dkb/s 下行 %6dkb/s 平均%6dkb/s\r", up, avgUp, down, avgDown)
}
func (p *speedPrinter) run() {
	var r *record
	var up, avgUp, down, avgDown int64
	for {
		select {
		case r = <-p.sendSpeedChan:
			up = r.Speed
			avgUp = r.AvgSpeed
			printSpeed(up, avgUp, down, avgDown)
		case r= <- p.recvSpeedChan:
			down = r.Speed
			avgDown = r.AvgSpeed
			printSpeed(up, avgUp, down, avgDown)
		}
	}
}