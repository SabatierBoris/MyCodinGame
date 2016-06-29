package main

import (
	"fmt"
	"math"
	"os"
	"sync"
)

const (
	NbTurn      = 200
	XShift      = 700
	YShift      = 700
	Xsize       = 16000
	Ysize       = 9000
	ReleaseDist = 1600
	StunDist    = 1760
	BustMaxDist = 1760
	BustMinDist = 900
	GhostSpeed  = 400
	BusterSpeed = 800
	NbPaths     = 6
	NbShift     = 3
)

//=============================================================================
//= INPUT LINE ================================================================
//=============================================================================
type InputLine struct {
	EntityId   int
	X          int
	Y          int
	EntityType int
	State      int
	Value      int
}

func (i InputLine) String() string {
	return fmt.Sprintf("Id:%d X:%d Y:%d Type:%d State:%d Value:%d", i.EntityId, i.X, i.Y, i.EntityType, i.State, i.Value)
}

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

//=============================================================================
//= PATH ======================================================================
//=============================================================================
type Path struct {
	currentIndex int
	list         []*Point
}

func (p *Path) Push(pt *Point) {
	p.list = append(p.list, pt)
}

func (p Path) GetCurrentPoint() *Point {
	if p.currentIndex < len(p.list) {
		return p.list[p.currentIndex]
	}
	return nil
}

func (p *Path) Next() {
	p.currentIndex++
}

func (p *Path) Reset() {
	p.currentIndex = 0
}

func (p Path) String() string {
	s := fmt.Sprintf("Current : %d", p.currentIndex)
	for i := 0; i < len(p.list); i++ {
		s = fmt.Sprintf("%s - %s", s, p.list[i])
	}
	return s
}

//=============================================================================
//= PATHS =====================================================================
//=============================================================================

func PathsService(in chan<- *Path, out <-chan *Path) {

}

type Paths struct {
	list []*Path
}

func (ps *Paths) Push(p *Path) {
	//fmt.Fprintf(os.Stderr, "Push path %s\n", p)
	ps.list = append(ps.list, p)
}

func (ps *Paths) Pop() *Path {
	p := ps.list[0]
	ps.list = ps.list[1:]
	//fmt.Fprintf(os.Stderr, "Pop path %s\n", p)
	return p
}

//=============================================================================
//= AGENT =====================================================================
//=============================================================================

type Agent struct {
	Id           int
	teamId       int
	quit         chan bool
	order        chan string
	data         chan InputLine
	endDigest    chan bool
	prepareOrder chan bool
	pos          Point
	paths        chan *Path
	currentPath  *Path
	reload       int
	orderSet     bool
	state        int
	value        int
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
			if data.EntityType == a.teamId && data.EntityId == a.Id {
				a.pos.X, a.pos.Y = data.X, data.Y
				a.state = data.State
				a.value = data.Value
				break
			}
			//fmt.Fprintf(os.Stderr, "Agent %d received : %s\n", a.Id, data)
			//TODO
		case <-a.endDigest:
			fmt.Fprintf(os.Stderr, "Agent %d End of digest\n", a.Id)
			//If I'm STUNED
			if a.state == 2 {
				a.order <- fmt.Sprintf("MOVE %s STUNED", a.pos)
				a.orderSet = true
				break
			}
			//TODO PREPARE ORDER FOR ATTACK (STUN or BUST)
			//TODO ASK HELP
			//TODO Tell I can HELP
		case <-a.prepareOrder:
			if a.orderSet == false {
				//Don't have order yet, Help someone
				//TODO
			}
			if a.orderSet == false {
				//Don't have order yet, use paths
				var p *Point
				p = nil
				for p == nil {
					if a.currentPath == nil {
						fmt.Fprintf(os.Stderr, "Agent %d get a new path\n", a.Id)
						a.currentPath = <-a.paths
						a.currentPath.Reset()
						fmt.Fprintf(os.Stderr, "Agent %d path %s\n", a.Id, a.currentPath)
					}
					p = a.currentPath.GetCurrentPoint()
					if p == nil {
						a.currentPath.Reset()
						a.paths <- a.currentPath
						a.currentPath = nil
					}
					if p != nil && p.GetDistanceTo(a.pos) < 100 {
						a.currentPath.Next()
						p = nil
					}
				}
				a.order <- fmt.Sprintf("MOVE %s %d", p, a.Id)
				a.orderSet = true
			}

			a.orderSet = false
			if a.reload > 0 {
				a.reload--
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
	input := InputLine{entityId, x, y, entityType, state, value}
	a.data <- input
}

func (a *Agent) EndDigestData() {
	a.endDigest <- true
}

func (a *Agent) PrepareOrder() {
	a.prepareOrder <- true
}

func MakeAgent(index, teamId, teamSize, nbGhost int, paths chan *Path) *Agent {
	agent := &Agent{index + (teamSize * teamId), teamId, make(chan bool), make(chan string, 1), make(chan InputLine), make(chan bool), make(chan bool), Point{0, 0}, paths, nil, 0, false, 0, 0}
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

	paths := make(chan *Path, NbPaths)

	for j := 0; j < NbShift; j++ {
		for i := j; i < NbPaths; i += NbShift {
			p := &Path{0, make([]*Point, 0)}
			p.Push(&Point{XShift + (i * ((Xsize - (2 * XShift)) / (NbPaths - 1))), Ysize - (YShift + (i * ((Ysize - (2 * YShift)) / (NbPaths - 1))))})
			if myTeamId == 0 {
				p.Push(&Point{Xsize - (4 * XShift), Ysize - (4 * YShift)})
				//TODO TRY p.Push(&Point{XShift + (i * ((Xsize - (2 * XShift)) / (NbPaths - 1))), Ysize - (YShift + (i * ((Ysize - (2 * YShift)) / (NbPaths - 1))))})
				p.Push(&Point{2 * XShift, 2 * YShift})
			} else {
				p.Push(&Point{2 * XShift, 2 * YShift})
				//TODO TRY p.Push(&Point{XShift + (i * ((Xsize - (2 * XShift)) / (NbPaths - 1))), Ysize - (YShift + (i * ((Ysize - (2 * YShift)) / (NbPaths - 1))))})
				p.Push(&Point{Xsize - (4 * XShift), Ysize - (4 * YShift)})
			}
			fmt.Fprintf(os.Stderr, "%d - Path : %s\n", myTeamId, p)
			paths <- p
		}
	}

	terminated.Add(bustersPerPlayer)
	for i := 0; i < bustersPerPlayer; i++ {
		agent := MakeAgent(i, myTeamId, bustersPerPlayer, ghostCount, paths)
		go agent.Run(&terminated)
		agents = append(agents, agent)
	}

	for current_turn := 0; current_turn < NbTurn; current_turn++ {
		fmt.Fprintf(os.Stderr, "Turn : %d\n", current_turn)
		var entities int
		fmt.Scan(&entities)
		for i := 0; i < entities; i++ {
			var entityId, x, y, entityType, state, value int
			// entityId: buster id or ghost id
			// y: position of this buster / ghost
			// entityType: the team id if it is a buster, -1 if it is a ghost.
			// state: For busters: 0=idle, 1=carrying a ghost.
			// value: For busters: Ghost id being carried. For ghosts: number of busters attempting to trap this ghost.
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
