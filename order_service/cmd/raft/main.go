package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/JuanGQCadavid/ds-practice-2025/order_service/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type StatesOfDemocracy string

const (
	Submission  StatesOfDemocracy = "(Follower) Well.. it could be worst" // Follower
	ALaLanterne StatesOfDemocracy = "(Candidate) À la lanterne!!"         // Candidate
	Reich       StatesOfDemocracy = "(Leader) I'm the Reich!"             // Leader
)

type MessageFromThePlebeTypes string
type MessageFromThePlebe struct {
	MesssageType MessageFromThePlebeTypes
	Term         int
}

const (
	// Messages
	CoupDeAaaahType           MessageFromThePlebeTypes = "Coup d'état!!"
	TheSonOfBitchIsStillAlive MessageFromThePlebeTypes = "Damn it! Not yet"

	// To hear about the other
	PROTOCOL = "tcp"
)

// -----------------------------
//
// SERVER
//
// -----------------------------

type Server struct {
	pb.UnimplementedConsensusServer

	messages         chan MessageFromThePlebe
	democracyUpdates <-chan StatesOfDemocracy

	latestStateOfDemocracy StatesOfDemocracy
}

var (
	// gRPC
	listener             net.Listener
	messagesFromThePlebe chan MessageFromThePlebe

	// Raft
	heartBeat = time.Duration(2 * time.Second)
	miniKafka *Kafki

	// Docker
	DEFAULT_PORT = os.Getenv("PORT")

	thePlebeAddress = []string{
		"localhost:50052",
		"localhost:50053",
		"localhost:50054",
		// "localhost:50055",
	}

	thePlebe []pb.ConsensusClient
)

func (srv *Server) InitServer() *Server {
	go func() {
		for msg := range srv.democracyUpdates {
			log.Println("New state saved on server", msg)
			srv.latestStateOfDemocracy = msg
		}
		log.Println("Server stop listening democracy updates")
	}()

	return srv
}

func (srv *Server) CoupDeAaaah(ctx context.Context, rq *pb.Empty) (*pb.CoupDEtatResponse, error) {
	log.Println("Oh shit! This is getting bananas!")

	srv.messages <- MessageFromThePlebe{
		MesssageType: CoupDeAaaahType,
	}

	if srv.latestStateOfDemocracy == Submission {
		return &pb.CoupDEtatResponse{
			Oks: true,
		}, nil
	}

	return &pb.CoupDEtatResponse{
		Oks: false,
	}, nil

}

func (srv *Server) YeahImStillAliveBitch(ctx context.Context, rq *pb.Empty) (*pb.Empty, error) {
	log.Println("Yet...")

	srv.messages <- MessageFromThePlebe{
		MesssageType: TheSonOfBitchIsStillAlive,
		Term:         int(rq.Term),
	}

	return &pb.Empty{}, nil
}

// -----------------------------
//
// # Inner Kafka
//
// -----------------------------

type Kafki struct {
	mu     sync.Mutex
	subs   []chan StatesOfDemocracy
	quit   chan struct{}
	closed bool
}

func NewAgent() *Kafki {
	return &Kafki{
		subs: make([]chan StatesOfDemocracy, 0),
		quit: make(chan struct{}),
	}
}

func (b *Kafki) Publish(state StatesOfDemocracy) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.closed {
		return
	}
	// THis could be abottle neck
	for _, ch := range b.subs {
		ch <- state
	}
}

func (b *Kafki) Subscribe() <-chan StatesOfDemocracy {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.closed {
		return nil
	}

	ch := make(chan StatesOfDemocracy, 100)
	b.subs = append(b.subs, ch)
	return ch
}

// -----------------------------
//
// # Main thread
//
// -----------------------------
func init() {
	if DEFAULT_PORT == "" {
		log.Fatalln("Missing port babe")
	}

	//Removing my self from the list // REMOVE ME
	newList := make([]string, 0)
	for _, addr := range thePlebeAddress {
		if !strings.Contains(addr, DEFAULT_PORT) {
			newList = append(newList, addr)
		}
	}
	thePlebeAddress = newList
	// DELETE ME

	var err error
	listener, err = net.Listen(PROTOCOL, fmt.Sprintf(":%s", DEFAULT_PORT))

	if err != nil {
		log.Panic("Unable to start listener in port ", DEFAULT_PORT, PROTOCOL, " err: ", err.Error())
	}
	messagesFromThePlebe = make(chan MessageFromThePlebe, 100)

	// Creating mini kafka
	miniKafka = NewAgent()

	// Creating connection with the Plebe
	thePlebe = make([]pb.ConsensusClient, len(thePlebeAddress))

	// I will try to create the connections to then, but not to stuck if they fail
	for pos, address := range thePlebeAddress {
		thePlebe[pos] = nil
		go func(pos int, dns string) {
			for {
				log.Println("knock knock!, Do you have time to talk about Jesus?", dns)
				conn, err := grpc.NewClient(dns, grpc.WithTransportCredentials(insecure.NewCredentials()))
				if err == nil {
					log.Println("Connection good with ", address)
					thePlebe[pos] = pb.NewConsensusClient(conn)
					break
				} else {
					log.Println("I will be back...")
					time.Sleep(1 * time.Second)
				}
			}
		}(pos, address)
	}

}

func startServer(thePlebeSpoke chan MessageFromThePlebe, democracyUpdates <-chan StatesOfDemocracy) {

	grpcServer := grpc.NewServer()

	server := &Server{
		messages:               thePlebeSpoke,
		democracyUpdates:       democracyUpdates,
		latestStateOfDemocracy: Submission,
	}
	server.InitServer()

	pb.RegisterConsensusServer(grpcServer, server)

	log.Println("Listening to all the madafuckers")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

func I_am_not_dead_yet(term int) {
	for _, plebe := range thePlebe {
		if plebe != nil {
			go func(plebe pb.ConsensusClient) {
				ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
				defer cancel()
				plebe.YeahImStillAliveBitch(ctx, &pb.Empty{
					Term: int64(term),
				})
			}(plebe)
		}
	}
}

func I_Have_A_Dream() StatesOfDemocracy {
	var (
		votes        int = 0
		wg           sync.WaitGroup
		voteMux      = sync.Mutex{}
		activePeople = len(thePlebe)
	)
	log.Println("I have a dream...")

	votes += 1 // Ofc! I'm the best option for this matorral.

	for _, plebe := range thePlebe {
		if plebe != nil {
			wg.Add(1)
			go func(plebe pb.ConsensusClient) {
				defer wg.Done()

				ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
				defer cancel()
				resp, err := plebe.CoupDeAaaah(ctx, &pb.Empty{})
				if err == nil {
					if resp.Oks {
						// Critical area
						voteMux.Lock()
						votes += 1
						voteMux.Unlock()
					}
				}
			}(plebe)
		} else {
			activePeople -= 1
		}

	}

	wg.Wait()

	log.Println("Votes: ", votes, "/", activePeople)
	if votes >= (activePeople / 2) {
		log.Println("The People's Will Has Spoken! A Pope Rises!")
		return Reich
	}

	return Submission
}

func main() {
	var (
		// Democracy
		actualState StatesOfDemocracy = Submission
		term        int

		// Election concerns
		baseElectionTimeOutMs = 150
		upToElectionTimeoutMs = 300
		electionTimeoutRandMs = baseElectionTimeOutMs + rand.Intn(upToElectionTimeoutMs-baseElectionTimeOutMs)
		electionTimeout       = time.Millisecond * time.Duration(electionTimeoutRandMs)

		iAmStillAliveDuration = 200 * time.Millisecond
	)

	go func() {
		serverDemocracyUpdates := miniKafka.Subscribe()
		startServer(messagesFromThePlebe, serverDemocracyUpdates)
	}()

	heartbit := time.NewTicker(heartBeat)

	electionTicker := time.NewTicker(electionTimeout)
	electionTicker.Stop()

	iAmStillAliveTimer := time.NewTicker(iAmStillAliveDuration)

	for {
		log.Println("-------------------")
		log.Println("I'am ", actualState, "Term", term)
		log.Println("-------------------")
		select {
		// Send life signals
		case <-iAmStillAliveTimer.C:
			if actualState == Reich {
				I_am_not_dead_yet(term)
			}

		// Gosipies
		case msg := <-messagesFromThePlebe:
			electionTicker.Stop()

			if msg.MesssageType == TheSonOfBitchIsStillAlive {
				log.Println("Damn it! He is alive...")

				if msg.Term > term {
					log.Println("Well, It seems I'm not longer the leader :'( ")
					term = msg.Term
					actualState = Submission
					miniKafka.Publish(Submission)
				}
			} else if msg.MesssageType == CoupDeAaaahType {
				log.Println("Someone start the democrazy! Jaaa, crazy...")
				actualState = Submission
				miniKafka.Publish(Submission)
			}

			heartbit.Reset(heartBeat)

		// My timers
		case time := <-heartbit.C:
			if actualState == Reich {
				continue
			}

			log.Println("Heartbit!!", time)
			log.Println("Lets prepare for the Independecy!!!")

			electionTicker.Reset(electionTimeout)
			heartbit.Reset(heartBeat)

		case time := <-electionTicker.C:
			log.Println("Today, we will write down a new chapter on the history fo this great nation, ", time)

			// No need to go back to election
			electionTicker.Stop()
			heartbit.Reset(heartBeat)

			actualState = I_Have_A_Dream()

			if actualState == Reich {
				term += 1
			}

			miniKafka.Publish(Submission)
		}
	}

}
