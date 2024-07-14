package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net"
	"time"

	"github.com/Jawadh-Salih/go-redis/client"
)

const defaultListenAddr = ":5001"

type Config struct {
	ListenAddress string
}

type Message struct {
	data []byte
	peer *Peer
}
type Server struct {
	Config
	peers     map[*Peer]bool
	ln        net.Listener
	addPeerch chan *Peer
	quitCh    chan struct{}
	msgCh     chan Message

	kv *KeyVal
}

func NewServer(cfg Config) *Server {
	if len(cfg.ListenAddress) == 0 {
		cfg.ListenAddress = defaultListenAddr
	}
	return &Server{
		Config:    cfg,
		peers:     make(map[*Peer]bool),
		addPeerch: make(chan *Peer),
		quitCh:    make(chan struct{}),
		msgCh:     make(chan Message),
		kv:        NewKeyVal(),
	}
}

func (s *Server) loop() {
	for {
		select {
		case rawMsg := <-s.msgCh:
			if err := s.handleRawMessage(rawMsg); err != nil {
				slog.Error("raw message error", "err", err)
			}
		case <-s.quitCh:
			return
		case peer := <-s.addPeerch:
			s.peers[peer] = true
		}
	}
}

func (s *Server) set(key string, val string) error {
	return nil
}

func (s *Server) handleRawMessage(msg Message) error {
	cmd, err := parseCommand(string(msg.data))
	if err != nil {
		return err
	}

	switch cm := cmd.(type) {
	case SetCommand:
		slog.Info("Set a key to the hash table", "command", cm)
		return s.kv.Set(cm.key, cm.val)

	case GetCommand:
		val, ok := s.kv.Get(cm.key)
		if !ok {
			return fmt.Errorf("key not found - %s", string(val))
		}

		_, err := msg.peer.Send(val)
		if err != nil {
			slog.Error("Peer send error", "err", err)
		}
	}

	return nil
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.ListenAddress)
	if err != nil {
		return err
	}

	s.ln = ln

	go s.loop()

	slog.Info("server running", "listenaddress", s.ListenAddress)

	return s.acceptLoop()
}

func (s *Server) acceptLoop() error {
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			slog.Error("Accept error", "err", err)
			continue
		}

		go s.handleConn(conn)
	}
}

func (s *Server) handleConn(conn net.Conn) {
	peer := NewPeer(conn, s.msgCh)

	s.addPeerch <- peer

	// > telnet localhost 5001
	// creattes a new peer to the server and connects
	slog.Info("new peer connected", "remoteaddress", conn.RemoteAddr())
	if err := peer.readLoop(); err != nil {
		slog.Error("peer read error", "err", err, "remoteAddr", conn.RemoteAddr())
	}
}
func main() {
	server := NewServer(Config{})
	go func() {
		log.Fatal(server.Start())

	}()

	time.Sleep(time.Second) // because the server needs to fully functional

	client := client.New("localhost:5001")
	for i := 0; i < 10; i++ {
		if err := client.Set(context.TODO(), fmt.Sprintf("foo-%d", i), fmt.Sprintf("bar-%d", i)); err != nil {
			log.Fatal(err)
		}

		time.Sleep(time.Second)

		val, err := client.Get(context.TODO(), fmt.Sprintf("foo-%d", i))
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Got thew valueu", val)
	}

	// fmt.Println(server.kv.data)
	// time.Sleep(2 * time.Second)

	// we are blocking herer. so the program does not exist
	// select {}
}
