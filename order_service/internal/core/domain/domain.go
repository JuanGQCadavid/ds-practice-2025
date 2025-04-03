package domain

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
)
