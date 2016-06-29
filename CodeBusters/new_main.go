package main

import (
	"fmt"
	"math"
	"os"
	"sync"
)

const (
	NbTurn       = 200
	XShift       = 700
	YShift       = 700
	Xsize        = 16000
	Ysize        = 9000
	ReleaseDist  = 1600
	StunDist     = 1760
	BustMaxDist  = 1760
	BustMinDist  = 900
	GhostSpeed   = 400
	BusterSpeed  = 800
	NbCheckpoint = 6
	NbShift      = 3
)

//=============================================================================
//= POINT =====================================================================
//=============================================================================
type Point struct {
	X int
	Y int
}

func (p Point) String() string {
	return fmt.Sprintf("%d %d", p.X, p.Y)
}

func (p Point) GetDistanceTo(t Point) float64 {
	dx := (float64)(p.X - t.X)
	dy := (float64)(p.Y - t.Y)
	return math.Hypot(dx, dy)
}

func (p Point) GetPositionAwaysFrom(t Point, distance float64) Point {
	//TODO
	return p
}

var Bases = []Point{Point{0, 0}, Point{Xsize, Ysize}}

//TODO Define TEAM BASES

//=============================================================================
//= PATH ======================================================================
//=============================================================================
//TODO

//=============================================================================
//= PATHS =====================================================================
//=============================================================================
//TODO

//=============================================================================
//= AGENT =====================================================================
//=============================================================================

type Agent struct {
	Id           int
	teamId       int
	quit         chan bool
	order        chan string
	data         chan bool //TODO Change data chan type
	endDigest    chan bool
	prepareOrder chan bool
}

func (a *Agent) Run(terminated *sync.WaitGroup) {
	defer terminated.Done()

	fmt.Fprintf(os.Stderr, "Agent %d running\n", a.Id)
	for {
		select {
		case <-a.quit:
			fmt.Fprintf(os.Stderr, "Agent %d terminated\n", a.Id)
			return
		case data := <-a.data:
			fmt.Fprintf(os.Stderr, "Agent %d received : %s\n", a.Id, data)
			//TODO
		case <-a.endDigest:
			fmt.Fprintf(os.Stderr, "Agent %d End of digest\n", a.Id)
			//TODO PREPARE ORDER FOR ATTACK (STUN or BUST)
			//TODO ASK HELP
			//TODO Tell I can HELP
		case <-a.prepareOrder:
			fmt.Fprintf(os.Stderr, "Agent %d Prepare order\n", a.Id)
			//TODO If I don't have order already
			//TODO Help someone
			//TODO Move somewhere if no one need help
			if a.teamId == 0 {
				a.order <- fmt.Sprintf("MOVE %s %d", Bases[1], a.Id)
			} else {
				a.order <- fmt.Sprintf("MOVE %s %d", Bases[0], a.Id)
			}
		}
	}
}

func (a *Agent) Terminate() {
	//fmt.Fprintf(os.Stderr, "Send terminate to %d\n", a.Id)
	a.quit <- true
}

func (a *Agent) GetOrder() string {
	//fmt.Fprintf(os.Stderr, "Wait order of %d\n", a.Id)
	return <-a.order
}

func (a *Agent) DigestData(entityId, x, y, entityType, state, value int) {
	//TODO Do something
	a.data <- true
}

func (a *Agent) EndDigestData() {
	a.endDigest <- true
}

func (a *Agent) PrepareOrder() {
	a.prepareOrder <- true
}

func MakeAgent(index, teamId, teamSize, nbGhost int) *Agent {
	agent := &Agent{index + (teamSize * teamId), teamId, make(chan bool), make(chan string), make(chan bool), make(chan bool), make(chan bool)}
	return agent
}

//=============================================================================
//= MAIN ======================================================================
//=============================================================================

func main() {
	var agents []*Agent
	var terminated sync.WaitGroup

	var bustersPerPlayer int
	var ghostCount int
	var myTeamId int
	fmt.Scan(&bustersPerPlayer)
	fmt.Scan(&ghostCount)
	fmt.Scan(&myTeamId)

	terminated.Add(bustersPerPlayer)
	for i := 0; i < bustersPerPlayer; i++ {
		agent := MakeAgent(i, myTeamId, bustersPerPlayer, ghostCount)
		go agent.Run(&terminated)
		agents = append(agents, agent)
	}

	for current_turn := 0; current_turn < NbTurn; current_turn++ {
		fmt.Fprintf(os.Stderr, "Turn : %d\n", current_turn)
		var entities int
		fmt.Scan(&entities)
		for i := 0; i < entities; i++ {
			// entityId: buster id or ghost id
			// y: position of this buster / ghost
			// entityType: the team id if it is a buster, -1 if it is a ghost.
			// state: For busters: 0=idle, 1=carrying a ghost.
			// value: For busters: Ghost id being carried. For ghosts: number of busters attempting to trap this ghost.
			var entityId, x, y, entityType, state, value int
			fmt.Scan(&entityId, &x, &y, &entityType, &state, &value)
			for j := 0; j < bustersPerPlayer; j++ {
				agents[j].DigestData(entityId, x, y, entityType, state, value)
			}
		}
		for i := 0; i < bustersPerPlayer; i++ {
			agents[i].EndDigestData()
		}
		for i := 0; i < bustersPerPlayer; i++ {
			agents[i].PrepareOrder()
		}
		for i := 0; i < bustersPerPlayer; i++ {
			fmt.Printf("%s\n", agents[i].GetOrder())
		}
	}

	for i := 0; i < bustersPerPlayer; i++ {
		agents[i].Terminate()
	}

	terminated.Wait()
	fmt.Fprintf(os.Stderr, "Main terminated\n")
}
