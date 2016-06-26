package main

import (
	"fmt"
	"math"
	"os"
	"sort"
)

type Point struct {
	X int
	Y int
}

func (p Point) String() string {
	return fmt.Sprintf("%d %d", p.X, p.Y)
}

func (p *Point) Update(x, y int) {
	p.X = x
	p.Y = y
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

type Checkpoints struct {
	list []*Point
}

func (c *Checkpoints) Push(p *Point) {
	c.list = append(c.list, p)
}

func (c *Checkpoints) Pop() *Point {
	p := c.list[0]
	c.list = c.list[1:]
	return p
}

type Buster struct {
	Id      int
	Pos     Point
	Value   int
	State   int
	Target  *Point
	Reload  int
	Visible bool
}

func (b *Buster) Update(x, y, state, value int) {
	b.Pos.Update(x, y)
	b.State = state
	b.Value = value
	if b.Reload > 0 {
		b.Reload--
	}
}

func (b Buster) String() string {
	return fmt.Sprintf("%d : %s - State:%d - Value:%d - Reload:%d", b.Id, b.Pos, b.State, b.Value, b.Reload)
}

type Ghost struct {
	Id     int
	Pos    Point
	Value  int
	IsSeen bool
}

func (g *Ghost) Update(x, y, value int, seen bool) {
	g.Pos.Update(x, y)
	g.Value = value
	g.IsSeen = seen
}

func (g Ghost) String() string {
	if g.IsSeen {
		return fmt.Sprintf("%d at %s (%s)", g.Id, g.Pos, "Visible")
	} else {
		return fmt.Sprintf("%d at %s (%s)", g.Id, g.Pos, "Hidden")
	}

}

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

type Team struct {
	TeamId      int
	Size        int
	Members     []Buster
	Base        Point
	KnownGhosts []Ghost
	Opponents   []Buster
	checkpoints Checkpoints
}

func (t *Team) Update() {
	// entities: the number of busters and ghosts visible to you
	var entities int
	fmt.Scan(&entities)

	for index, _ := range t.KnownGhosts {
		t.KnownGhosts[index].IsSeen = false
	}

	for index, _ := range t.Opponents {
		t.Opponents[index].Visible = false
	}

	for i := 0; i < entities; i++ {
		// entityId: buster id or ghost id
		// y: position of this buster / ghost
		// entityType: the team id if it is a buster, -1 if it is a ghost.
		// state: For busters: 0=idle, 1=carrying a ghost.
		// value: For busters: Ghost id being carried. For ghosts: number of busters attempting to trap this ghost.
		var entityId, x, y, entityType, state, value int
		fmt.Scan(&entityId, &x, &y, &entityType, &state, &value)

		if entityType == -1 { //GHOST
			t.UpdateOrCreateGhost(entityId, x, y, value)
		} else if entityType == t.TeamId { //My Team
			t.UpdateMember(entityId, x, y, state, value)
			if value != -1 {
				t.RemoveGhost(value)
			}
		} else { //Opponents
			t.UpdateOpponent(entityId, x, y, state, value)
			if value != -1 {
				t.RemoveGhost(value)
			}
		}
	}
}

func (t *Team) RemoveGhost(value int) {
	for index, ghost := range t.KnownGhosts {
		if ghost.Id == value {
			fmt.Fprintf(os.Stderr, "Delete ghost %s\n", ghost)
			t.KnownGhosts = append(t.KnownGhosts[:index], t.KnownGhosts[index+1:]...)
			return
		}
	}
}

func (t *Team) UpdateOrCreateGhost(entityId, x, y, value int) {
	for index, ghost := range t.KnownGhosts {
		if ghost.Id == entityId { //If ghost is already known
			//Update it
			t.KnownGhosts[index].Update(x, y, value, true)
			//fmt.Fprintf(os.Stderr, "Update ghost %s\n", t.KnownGhosts[index])
			return
		}
	}
	//If ghost isn't found
	//Create it
	ghost := Ghost{entityId, Point{x, y}, value, true}
	//fmt.Fprintf(os.Stderr, "Add new ghost %s\n", ghost)
	t.KnownGhosts = append(t.KnownGhosts, ghost)

}

func (t *Team) UpdateMember(entityId, x, y, state, value int) {
	index := entityId
	if t.TeamId == 1 {
		index -= t.Size
	}
	t.Members[index].Update(x, y, state, value)
	//fmt.Fprintf(os.Stderr, " Update member %d  : %s\n", index, t.Members[index])
}

func (t *Team) UpdateOpponent(entityId, x, y, state, value int) {
	index := entityId
	if t.TeamId == 0 {
		index -= t.Size
	}
	t.Opponents[index].Update(x, y, state, value)
	t.Opponents[index].Visible = true
	fmt.Fprintf(os.Stderr, " Update opponent %d  : %s\n", index, t.Opponents[index])
}

func (t Team) GetOrderedGhostByDistanceOf(p Point) []*Ghost {
	buf := Ghosts{make([]*Ghost, len(t.KnownGhosts)), p}
	for index, _ := range t.KnownGhosts {
		buf.List[index] = &t.KnownGhosts[index]
	}
	sort.Sort(buf)
	//fmt.Fprintf(os.Stderr, "%s\n", buf)
	return buf.List
}

func (t Team) GetNearestFreeMemberOf(p Point) *Buster {
	var nearestMember *Buster
	var nearestDist float64
	nearestMember = nil
	for index, member := range t.Members {
		if member.State == 0 {
			dist := member.Pos.GetDistanceTo(p)
			if nearestMember == nil || dist < nearestDist {
				nearestMember = &t.Members[index]
				nearestDist = dist
			}
		}
	}
	return nearestMember
}

func (t Team) GetReachableOpponent(b Buster) *Buster {
	for index, _ := range t.Opponents {
		if t.Opponents[index].Visible && t.Opponents[index].Pos.GetDistanceTo(b.Pos) < 1760 {
			return &t.Opponents[index]
		}
	}
	return nil
}

func (t *Team) DisplayOrders() {
	for i := 0; i < t.Size; i++ {
		if t.Members[i].State == 1 {
			//If the members have a ghost
			dist := t.Members[i].Pos.GetDistanceTo(t.Base)
			fmt.Fprintf(os.Stderr, "%d have a ghost (Distance to %s : %f - %t)\n", i, t.Base, dist, (dist > 1600))
			if dist > 1600 {
				fmt.Printf("MOVE %s\n", t.Base)
			} else {
				fmt.Printf("RELEASE\n")
			}
		} else {
			if t.Members[i].Reload == 0 { //Can shoot
				nearestOpponent := t.GetReachableOpponent(t.Members[i])
				if nearestOpponent != nil {
					fmt.Printf("STUN %d\n", nearestOpponent.Id)
					t.Members[i].Reload = 20
					continue
				}
			}

			//fmt.Fprintf(os.Stderr, "%d haven't a ghost\n", i)
			ghosts := t.GetOrderedGhostByDistanceOf(t.Members[i].Pos)
			order := false
			for _, ghost := range ghosts {
				nearestMember := t.GetNearestFreeMemberOf(ghost.Pos)
				if nearestMember == &t.Members[i] {
					dist := t.Members[i].Pos.GetDistanceTo(ghost.Pos)
					fmt.Fprintf(os.Stderr, "%d target %s (dist:%f)\n", t.Members[i].Id, ghost, dist)
					if dist > 1760 {
						fmt.Printf("MOVE %s\n", ghost.Pos)
						order = true
						break
					} else if dist < 900 && ghost.IsSeen {
						if dist+400 > 900 {
							fmt.Printf("MOVE %s\n", t.Members[i].Pos)
						} else {
							//Move out
							targetPos := t.Members[i].Pos.GetPositionAwaysFrom(ghost.Pos, 900-(dist+400))
							fmt.Printf("MOVE %s\n", targetPos)
						}
						order = true
						break
					} else if ghost.IsSeen {
						fmt.Printf("BUST %d\n", ghost.Id)
						order = true
						break
					} else {
						t.RemoveGhost(ghost.Id)
					}
				}
			}
			if !order {
				if t.Members[i].Target == nil || t.Members[i].Target.GetDistanceTo(t.Members[i].Pos) < 100 {
					if t.Members[i].Target != nil {
						t.checkpoints.Push(t.Members[i].Target)
					}
					t.Members[i].Target = t.checkpoints.Pop()
				}
				fmt.Printf("MOVE %s\n", t.Members[i].Target)
			}
		}
	}
}

func CreateTeam(size int, id int) *Team {
	var t *Team
	t = nil
	switch id {
	case 0:
		t = &Team{id, size, make([]Buster, size), Point{0, 0}, make([]Ghost, 0), make([]Buster, size), Checkpoints{}}
		t.checkpoints.Push(&Point{1000, 8000})
		t.checkpoints.Push(&Point{5000, 8000})
		t.checkpoints.Push(&Point{9000, 8000})
		t.checkpoints.Push(&Point{13000, 8000})
		t.checkpoints.Push(&Point{15000, 6000})
		t.checkpoints.Push(&Point{15000, 2000})
		t.checkpoints.Push(&Point{3000, 8000})
		t.checkpoints.Push(&Point{7000, 8000})
		t.checkpoints.Push(&Point{11000, 8000})
		t.checkpoints.Push(&Point{15000, 8000})
		t.checkpoints.Push(&Point{15000, 4000})
		t.checkpoints.Push(&Point{15000, 1})
	case 1:
		t = &Team{id, size, make([]Buster, size), Point{16000, 9000}, make([]Ghost, 0), make([]Buster, size), Checkpoints{}}
		t.checkpoints.Push(&Point{15000, 1000})
		t.checkpoints.Push(&Point{11000, 1000})
		t.checkpoints.Push(&Point{7000, 1000})
		t.checkpoints.Push(&Point{3000, 1000})
		t.checkpoints.Push(&Point{1000, 2000})
		t.checkpoints.Push(&Point{1000, 4000})
		t.checkpoints.Push(&Point{1000, 6000})
		t.checkpoints.Push(&Point{1000, 8000})
		t.checkpoints.Push(&Point{13000, 1000})
		t.checkpoints.Push(&Point{9000, 1000})
		t.checkpoints.Push(&Point{5000, 1000})
		t.checkpoints.Push(&Point{1000, 1})
	}
	for index, _ := range t.Members {
		memberIndex := index
		if t.TeamId == 1 {
			memberIndex += t.Size
		}
		t.Members[index].Id = memberIndex
	}
	for index, _ := range t.Opponents {
		memberIndex := index
		if t.TeamId == 0 {
			memberIndex += t.Size
		}
		t.Opponents[index].Id = memberIndex
	}
	return t
}

func main() {
	var bustersPerPlayer int
	var ghostCount int
	var myTeamId int
	fmt.Scan(&bustersPerPlayer)
	fmt.Scan(&ghostCount)
	fmt.Scan(&myTeamId)

	MyTeam := CreateTeam(bustersPerPlayer, myTeamId)
	for {
		MyTeam.Update()
		MyTeam.DisplayOrders()
	}
}
