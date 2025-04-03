package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"sync"
	"time"

	"github.com/JuanGQCadavid/ds-practice-2025/order_service/internal/core"
	"github.com/JuanGQCadavid/ds-practice-2025/order_service/internal/core/domain"
	"github.com/JuanGQCadavid/ds-practice-2025/order_service/internal/repositories/queue"
	"github.com/JuanGQCadavid/ds-practice-2025/order_service/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	// To hear about the other
	PROTOCOL       = "tcp"
	QUEUE_ENV_NAME = "order_queue_dns"

	PLEBE_1_DNS_ENV_NAME = "order_service_1_dns"
	PLEBE_2_DNS_ENV_NAME = "order_service_2_dns"
)

var (
	// gRPC
	listener             net.Listener
	messagesFromThePlebe chan domain.MessageFromThePlebe

	// Raft
	heartBeat = time.Duration(2 * time.Second)
	miniKafka *core.Kafki

	// Puller
	repo *queue.QueueRepository

	// Docker
	DEFAULT_PORT    = os.Getenv("PORT")
	thePlebeAddress []string

	thePlebe []pb.ConsensusClient
)

// -----------------------------
//
// # Main thread
//
// -----------------------------
func init() {
	// Preparing the plebe

	plebe1, ok := os.LookupEnv(PLEBE_1_DNS_ENV_NAME)

	if !ok {
		log.Fatalln("Missing the plebe 1 dns")
	}

	plebe2, ok := os.LookupEnv(PLEBE_2_DNS_ENV_NAME)

	if !ok {
		log.Fatalln("Missing plebe 2 dns")
	}

	thePlebeAddress = []string{
		plebe1,
		plebe2,
	}

	if DEFAULT_PORT == "" {
		log.Fatalln("Missing port babe")
	}

	var err error
	listener, err = net.Listen(PROTOCOL, fmt.Sprintf(":%s", DEFAULT_PORT))

	if err != nil {
		log.Panic("Unable to start listener in port ", DEFAULT_PORT, PROTOCOL, " err: ", err.Error())
	}
	messagesFromThePlebe = make(chan domain.MessageFromThePlebe, 100)

	// Puller config
	dns, ok := os.LookupEnv(QUEUE_ENV_NAME)

	if !ok {
		log.Fatalln("Missing QUEUE dns")
	}

	repo = queue.NewQueueRepository(dns)

	// Creating mini kafka
	miniKafka = core.NewAgent()

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

func startServer(thePlebeSpoke chan domain.MessageFromThePlebe, democracyUpdates <-chan domain.StatesOfDemocracy) {
	var (
		grpcServer = grpc.NewServer()
		server     = core.NewServer(thePlebeSpoke, democracyUpdates, domain.Submission)
	)

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
				if _, err := plebe.YeahImStillAliveBitch(ctx, &pb.Empty{
					Term: int64(term),
				}); err != nil {
					log.Println("I cant not send a hear bit to ", plebe, " error: ", err.Error())
				}

			}(plebe)
		}
	}
}

func I_Have_A_Dream() domain.StatesOfDemocracy {
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
				} else {
					log.Println("I cant not send a votation process to ", plebe, " error: ", err.Error())
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
		return domain.Reich
	}

	return domain.Submission
}

func main() {
	var (
		// Democracy
		actualState domain.StatesOfDemocracy = domain.Submission
		prevState   domain.StatesOfDemocracy = domain.ALaLanterne // Just a trick for getting a log...

		term int

		// Election concerns
		baseElectionTimeOutMs = 200  //150
		upToElectionTimeoutMs = 1000 //300
		electionTimeoutRandMs = baseElectionTimeOutMs + rand.Intn(upToElectionTimeoutMs-baseElectionTimeOutMs)
		electionTimeout       = time.Millisecond * time.Duration(electionTimeoutRandMs)

		iAmStillAliveDuration = 200 * time.Millisecond
	)

	log.Println("Election MS: ", electionTimeoutRandMs)

	// gRPC
	go func() {
		serverDemocracyUpdates := miniKafka.Subscribe()
		startServer(messagesFromThePlebe, serverDemocracyUpdates)
	}()

	// The puller
	go func() {
		serverDemocracyUpdates := miniKafka.Subscribe()
		service := core.NewService(repo, serverDemocracyUpdates)
		service.Init()
		service.Listen()
	}()

	heartbit := time.NewTicker(heartBeat)

	electionTicker := time.NewTicker(electionTimeout)
	electionTicker.Stop()

	iAmStillAliveTimer := time.NewTicker(iAmStillAliveDuration)

	for {

		if actualState != prevState {
			log.Println("-------------------")
			log.Println("I'am ", actualState, "Term", term)
			log.Println("-------------------")
			prevState = actualState
		}

		select {
		// Send life signals
		case <-iAmStillAliveTimer.C:
			if actualState == domain.Reich {
				I_am_not_dead_yet(term)
			}

		// Gosipies
		case msg := <-messagesFromThePlebe:
			electionTicker.Stop()

			if msg.MesssageType == domain.TheSonOfBitchIsStillAlive {
				// log.Println("Damn it! He is alive...")

				if msg.Term > term {
					log.Println("Well, It seems I'm not longer the leader :'( ")
					term = msg.Term
					actualState = domain.Submission
					miniKafka.Publish(domain.Submission)
				}
			} else if msg.MesssageType == domain.CoupDeAaaahType {
				log.Println("Someone start the democrazy! Jaaa, crazy...")
				actualState = domain.Submission
				miniKafka.Publish(domain.Submission)
			}

			heartbit.Reset(heartBeat)

		// My timers
		case time := <-heartbit.C:
			// case <-heartbit.C:
			if actualState == domain.Reich {
				continue
			}

			log.Println("I did not hear from the leader' for a while', time for a change..", time)

			electionTicker.Reset(electionTimeout)
			heartbit.Reset(heartBeat)

		case time := <-electionTicker.C:
			log.Println("Today, we will write down a new chapter on the history fo this great nation, ", time)

			electionTicker.Stop()
			heartbit.Reset(heartBeat)

			actualState = I_Have_A_Dream()

			if actualState == domain.Reich {
				term += 1
			}

			miniKafka.Publish(domain.Reich)
		}
	}
}
