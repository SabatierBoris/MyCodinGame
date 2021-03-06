package main

import (
	"fmt"
	"math"
	"os"
	"sort"
	"sync"
)

const (
	NbTurn          = 200
	XShift          = 300
	YShift          = 300
	Xsize           = 16000
	Ysize           = 9000
	ReleaseDist     = 1600
	StunDist        = 1760
	BustMaxDist     = 1760
	BustMinDist     = 900
	GhostSpeed      = 400
	BusterSpeed     = 800
	NbPaths         = 6
	NbShift         = 3
	FogDistance     = 2200
	WeakGhost       = 5
	AverageGhost    = 15
	NbTurnHighGhost = 30
)

const (
	_                = iota
	HighHelpLevel    = iota
	AverageHelpLevel = iota
	LowHelpLevel     = iota
	UselessHelpLevel = iota
)

//=============================================================================
//= SHARE HELPER ==============================================================
//=============================================================================
type ShareHelper struct {
	mut     *sync.Mutex
	helpers []*Agent
}

func (s *ShareHelper) AddHelper(helper *Agent) {
	s.mut.Lock()
	defer s.mut.Unlock()
	s.helpers = append(s.helpers, helper)
}

func (s *ShareHelper) Reset() {
	fmt.Fprintf(os.Stderr, "HelperList : %s\n", s)
	s.mut.Lock()
	defer s.mut.Unlock()
	s.helpers = []*Agent{}
}

func (s *ShareHelper) String() string {
	s.mut.Lock()
	defer s.mut.Unlock()
	ret := ""
	for _, helper := range s.helpers {
		ret = fmt.Sprintf("%s - %s", ret, helper)
	}
	return ret
}

//=============================================================================
//= HELP REQUEST ==============================================================
//=============================================================================
type HelpRequest struct {
	Pos          Point
	Level        int
	Count        int
	RequestCount int
}

func (h *HelpRequest) IncreaseCounter() {
	h.Count++
}

func (h HelpRequest) String() string {
	return fmt.Sprintf("Pos:%s Level:%d Count:%d", h.Pos, h.Level, h.Count)
}

func (h HelpRequest) GetScore(target *Point) float64 {
	return h.Pos.GetDistanceTo(*target) * (float64)(h.Level+h.Count)
}

func (h HelpRequest) IsAnwsered() bool {
	return h.Count >= h.RequestCount
}

//=============================================================================
//= SHARE HELP REQUESTS =======================================================
//=============================================================================
type ShareHelpRequests struct {
	mut      *sync.Mutex
	requests []*HelpRequest
}

func (s *ShareHelpRequests) AddRequest(pos Point, level int, nb int) {
	s.mut.Lock()
	defer s.mut.Unlock()
	for index, _ := range s.requests {
		if s.requests[index].Pos.X == pos.X && s.requests[index].Pos.Y == pos.Y {
			//If the request is already ask, decrease urgence level
			s.requests[index].IncreaseCounter()
			return
		}
	}
	s.requests = append(s.requests, &HelpRequest{pos, level, 0, nb})
}

func (s *ShareHelpRequests) Reset() {
	fmt.Fprintf(os.Stderr, "HelpList : %s\n", s)
	s.mut.Lock()
	defer s.mut.Unlock()
	s.requests = []*HelpRequest{}
}

func (s *ShareHelpRequests) String() string {
	s.mut.Lock()
	defer s.mut.Unlock()
	ret := ""
	for _, request := range s.requests {
		ret = fmt.Sprintf("%s - %s", ret, request)
	}
	return ret
}

func (s *ShareHelpRequests) AnwserToHelp(current *Point) *Point {
	s.mut.Lock()
	defer s.mut.Unlock()

	var request *HelpRequest
	request = nil
	var request_score float64

	for index, _ := range s.requests {
		current_score := s.requests[index].GetScore(current)
		if !s.requests[index].IsAnwsered() {
			if request == nil || current_score < request_score {
				request = s.requests[index]
				request_score = current_score
			}
		}
	}
	if request != nil {
		request.IncreaseCounter()
		return &request.Pos
	}
	return nil
}

//=============================================================================
//= ORDER =====================================================================
//=============================================================================
type Order struct {
	Id    int
	Count int
}

func (o Order) String() string {
	return fmt.Sprintf("Id:%d Count:%d", o.Id, o.Count)
}

//=============================================================================
//= SHARE ORDERS ==============================================================
//=============================================================================
type ShareOrders struct {
	mut    *sync.Mutex
	orders []*Order
}

func (s *ShareOrders) MakeOrder(order int, limit int) bool {
	s.mut.Lock()
	defer s.mut.Unlock()
	for index, _ := range s.orders {
		if s.orders[index].Id == order {
			if s.orders[index].Count < limit {
				s.orders[index].Count++
				return true
			} else {
				return false
			}
		}
	}
	newOrder := &Order{order, 1}
	s.orders = append(s.orders, newOrder)
	return true
}

func (s *ShareOrders) GetOrderCount(order int) int {
	s.mut.Lock()
	defer s.mut.Unlock()
	for index, _ := range s.orders {
		if s.orders[index].Id == order {
			return s.orders[index].Count
		}
	}
	return 0
}

func (s *ShareOrders) Reset() {
	fmt.Fprintf(os.Stderr, "OrdersList : %s\n", s)
	s.mut.Lock()
	defer s.mut.Unlock()
	s.orders = []*Order{}
}

func (s *ShareOrders) String() string {
	s.mut.Lock()
	defer s.mut.Unlock()
	ret := ""
	for _, order := range s.orders {
		ret = fmt.Sprintf("%s - %s", ret, order)
	}
	return ret
}

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
//= Ghost =====================================================================
//=============================================================================
type Ghost struct {
	Id    int
	Pos   Point
	State int
	Value int
}

//=============================================================================
//= Ghosts ====================================================================
//=============================================================================

type Ghosts struct {
	List   []*Ghost
	Target Point
}

func (slice Ghosts) Len() int {
	return len(slice.List)
}

func (slice Ghosts) String() string {
	ret := fmt.Sprintf("Target : %s\n", slice.Target)
	for _, ghost := range slice.List {
		ret = fmt.Sprintf("%s : %s\n", ret, ghost)
	}
	return ret
}

func (slice Ghosts) Less(i, j int) bool {
	return slice.List[i].Pos.GetDistanceTo(slice.Target) < slice.List[j].Pos.GetDistanceTo(slice.Target)
}

func (slice Ghosts) Swap(i, j int) {
	slice.List[i], slice.List[j] = slice.List[j], slice.List[i]
}

//=============================================================================
//= Opponent ==================================================================
//=============================================================================
type Opponent struct {
	Id    int
	Pos   Point
	State int
	Value int
}

//=============================================================================
//= Opponents =================================================================
//=============================================================================

type Opponents struct {
	List   []*Opponent
	Target Point
}

func (slice Opponents) Len() int {
	return len(slice.List)
}

func (slice Opponents) String() string {
	ret := fmt.Sprintf("Target : %s\n", slice.Target)
	for _, opponent := range slice.List {
		ret = fmt.Sprintf("%s : %s\n", ret, opponent)
	}
	return ret
}

func (slice Opponents) Less(i, j int) bool {
	return slice.List[i].Pos.GetDistanceTo(slice.Target) < slice.List[j].Pos.GetDistanceTo(slice.Target)
}

func (slice Opponents) Swap(i, j int) {
	slice.List[i], slice.List[j] = slice.List[j], slice.List[i]
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
	datas        []InputLine
	ghosts       []Ghost
	opponents    []Opponent
	lastBust     int
	bustOrders   *ShareOrders
	stunOrders   *ShareOrders
	helpList     *ShareHelpRequests
	helperList   *ShareHelper
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
			} else {
				a.datas = append(a.datas, data)
			}
		case <-a.endDigest:
			//fmt.Fprintf(os.Stderr, "Agent %d End of digest\n", a.Id)
			// Analyse datas
			if a.orderSet == false {
				for _, data := range a.datas {
					if data.EntityType == -1 {
						//Ghost
						p := Point{data.X, data.Y}
						if p.GetDistanceTo(a.pos) <= FogDistance {
							//fmt.Fprintf(os.Stderr, "Ghost near %d\n", a.Id)
							a.ghosts = append(a.ghosts, Ghost{data.EntityId, p, data.State, data.Value})
						}
					} else if data.EntityType != a.teamId {
						//Opponent
						p := Point{data.X, data.Y}
						if p.GetDistanceTo(a.pos) <= FogDistance {
							//fmt.Fprintf(os.Stderr, "Opponent near %d\n", a.Id)
							a.opponents = append(a.opponents, Opponent{data.EntityId, p, data.State, data.Value})
						}
					} else {
						//Friend
						//TODO
					}
				}
			}

			//If I'm STUNED
			if a.state == 2 {
				a.order <- fmt.Sprintf("MOVE %s Zzzzzz", a.pos)
				a.orderSet = true
				nb := 0
				level := AverageHelpLevel
				for _, opponent := range a.opponents {
					if opponent.State != 2 {
						level = HighHelpLevel
					}
					if opponent.State != 2 || opponent.Value <= 2 {
						nb++
					}
				}
				if nb > 0 {
					a.AskHelp(a.pos, level, nb)
				}
				break
			}

			if a.orderSet == false && a.reload == 0 {
				//If I can shoot and enemy is near and enemy isn't stun
				sortedOpponents := Opponents{[]*Opponent{}, a.pos}
				for index, _ := range a.opponents {
					if a.pos.GetDistanceTo(a.opponents[index].Pos) <= StunDist {
						sortedOpponents.List = append(sortedOpponents.List, &a.opponents[index])
					}
				}
				sort.Sort(sortedOpponents)
				for _, opponent := range sortedOpponents.List {
					if opponent.State != 2 || opponent.Value <= 2 {
						a.AttackOpponent(*opponent)
						break
					}
				}
			}

			if a.orderSet == false {
				if a.state == 1 {
					//If I have a ghost
					a.GoHomeAndRelease()
					break
				}
			}

			if a.orderSet == false && a.reload == 0 {
				//If I can shoot and enemy is near and enemy isn't stun
				sortedOpponents := Opponents{make([]*Opponent, len(a.opponents)), a.pos}
				for index, _ := range a.opponents {
					sortedOpponents.List[index] = &a.opponents[index]
				}
				sort.Sort(sortedOpponents)
				for _, opponent := range sortedOpponents.List {
					if opponent.State != 2 || opponent.Value <= 2 {
						a.AttackOpponent(*opponent)
						break
					}
				}
			}

			if a.orderSet == false {
				// Get nearest ghost order by distance
				sortedGhosts := Ghosts{make([]*Ghost, len(a.ghosts)), a.pos}
				for index, _ := range a.ghosts {
					sortedGhosts.List[index] = &a.ghosts[index]
				}
				sort.Sort(sortedGhosts)
				for _, ghost := range sortedGhosts.List {
					if ghost.State <= WeakGhost {
						//Weak ghost
						if a.AttackGhost(*ghost, 1) {
							a.lastBust = 0
						}
						if ghost.State == 0 && ghost.Value > 1 {
							a.AskHelp(ghost.Pos, HighHelpLevel, ghost.Value)
						}
						break
					} else if ghost.State <= AverageGhost {
						if a.AttackGhost(*ghost, 2) {
							a.lastBust = 0
						}
						a.AskHelp(ghost.Pos, AverageHelpLevel, 1)
						break
					} else {
						if a.lastBust > NbTurnHighGhost {
							a.AttackGhost(*ghost, 3)
							if ghost.Value > 1 {
								a.AskHelp(ghost.Pos, AverageHelpLevel, 3)
							} else {
								a.AskHelp(ghost.Pos, LowHelpLevel, 3)
							}
							break
						}
					}
				}
			}
			if a.orderSet == false {
				a.ICanHelp()
			}
		case <-a.prepareOrder:
			if a.orderSet == false {
				//TODO I'm the best to answer to this help ?
				if pos := a.helpList.AnwserToHelp(&a.pos); pos != nil {
					a.order <- fmt.Sprintf("MOVE %s support", pos)
					a.orderSet = true
					a.lastBust = NbTurnHighGhost
				}
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
					if p != nil && p.GetDistanceTo(a.pos) < BusterSpeed {
						a.currentPath.Next()
						p = nil
					}
				}
				a.order <- fmt.Sprintf("MOVE %s help?", p)
				a.orderSet = true
			}

			a.orderSet = false
			if a.reload > 0 {
				a.reload--
			}
			a.datas = []InputLine{}
			a.ghosts = []Ghost{}
			a.opponents = []Opponent{}
			a.lastBust++
		}
	}
}

func (a *Agent) ICanHelp() {
	a.helperList.AddHelper(a)
}

func (a *Agent) AskHelp(pos Point, level int, nb int) {
	fmt.Fprintf(os.Stderr, "Agent %d need %d help at %s (%d)\n", a.Id, nb, pos, level)
	a.helpList.AddRequest(pos, level, nb)
}

func (a *Agent) GoHomeAndRelease() {
	dist := a.pos.GetDistanceTo(Bases[a.teamId])
	//fmt.Fprintf(os.Stderr, "%d have a ghost (Distance to %s : %f - %t)\n", i, t.Base, dist, (dist > ReleaseDist))
	if dist > ReleaseDist {
		a.order <- fmt.Sprintf("MOVE %s home", Bases[a.teamId])
		a.orderSet = true
	} else {
		a.order <- fmt.Sprintf("RELEASE")
		a.orderSet = true
	}
}

func (a *Agent) AttackOpponent(o Opponent) bool {
	if a.stunOrders.GetOrderCount(o.Id) == 0 {
		dist := a.pos.GetDistanceTo(o.Pos)
		if dist > StunDist {
			a.order <- fmt.Sprintf("MOVE %s fight", o.Pos)
			a.orderSet = true
		} else {
			if a.stunOrders.MakeOrder(o.Id, 1) {
				a.order <- fmt.Sprintf("STUN %d", o.Id)
				a.orderSet = true
				a.reload = 20
				return true
			}
		}
	}
	return false
}

func (a *Agent) AttackGhost(g Ghost, limit int) bool {
	if a.bustOrders.GetOrderCount(g.Id) < limit {
		dist := a.pos.GetDistanceTo(g.Pos)
		//fmt.Fprintf(os.Stderr, "%d target %s (dist:%f)\n", buster.Id, ghost, dist)
		if dist > BustMaxDist {
			a.order <- fmt.Sprintf("MOVE %s attack", g.Pos)
			//TODO Only if no one esle
			a.orderSet = true
		} else if dist < BustMinDist {
			//Move out
			targetPos := a.pos
			if dist+GhostSpeed < BustMinDist {
				targetPos = a.pos.GetPositionAwaysFrom(g.Pos, BustMinDist-(dist+GhostSpeed))
			}
			a.order <- fmt.Sprintf("MOVE %s getout", targetPos)
			//TODO Only if no one esle
			a.orderSet = true
		} else {
			if a.bustOrders.MakeOrder(g.Id, limit) {
				a.order <- fmt.Sprintf("BUST %d", g.Id)
				a.orderSet = true
				return true
			}
		}
	}
	return false
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

func MakeAgent(index, teamId, teamSize, nbGhost int, paths chan *Path, bustOrders *ShareOrders, stunOrders *ShareOrders, helpList *ShareHelpRequests, helperList *ShareHelper) *Agent {
	agent := &Agent{index + (teamSize * teamId), teamId, make(chan bool), make(chan string, 1), make(chan InputLine), make(chan bool), make(chan bool), Point{0, 0}, paths, nil, 0, false, 0, 0, []InputLine{}, []Ghost{}, []Opponent{}, 0, bustOrders, stunOrders, helpList, helperList}
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
				p.Push(&Point{XShift + (i * ((Xsize - (2 * XShift)) / (NbPaths - 1))), Ysize - (YShift + (i * ((Ysize - (2 * YShift)) / (NbPaths - 1))))})
				p.Push(&Point{2 * XShift, 2 * YShift})
			} else {
				p.Push(&Point{2 * XShift, 2 * YShift})
				p.Push(&Point{XShift + (i * ((Xsize - (2 * XShift)) / (NbPaths - 1))), Ysize - (YShift + (i * ((Ysize - (2 * YShift)) / (NbPaths - 1))))})
				p.Push(&Point{Xsize - (4 * XShift), Ysize - (4 * YShift)})
			}
			//fmt.Fprintf(os.Stderr, "%d - Path : %s\n", myTeamId, p)
			paths <- p
		}
	}

	bustOrders := &ShareOrders{&sync.Mutex{}, []*Order{}}
	stunOrders := &ShareOrders{&sync.Mutex{}, []*Order{}}
	helpList := &ShareHelpRequests{&sync.Mutex{}, []*HelpRequest{}}
	helperList := &ShareHelper{&sync.Mutex{}, []*Agent{}}

	terminated.Add(bustersPerPlayer)
	for i := 0; i < bustersPerPlayer; i++ {
		agent := MakeAgent(i, myTeamId, bustersPerPlayer, ghostCount, paths, bustOrders, stunOrders, helpList, helperList)
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
		bustOrders.Reset()
		stunOrders.Reset()
		helpList.Reset()
		helperList.Reset()
	}

	for i := 0; i < bustersPerPlayer; i++ {
		agents[i].Terminate()
	}

	terminated.Wait()
	fmt.Fprintf(os.Stderr, "Main terminated\n")
}
