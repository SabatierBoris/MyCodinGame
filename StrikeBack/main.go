package main

import (
	"fmt"
	"os"
)

type Point struct {
	X int
	Y int
}

type Pod struct {
	currentPosition  *Point
	previousPosition *Point
}

func (p *Pod) GetCurrentSpeed() (sX int, sY int) {
	sX = 0
	sY = 0
	return
}

type Journey struct {
	path     []*Point
	previous *Point
	complet  bool
}

func (j *Journey) GetCurrentTarget(x int, y int) *Point {
	for i := 0; i < len(j.path); i++ {
		if j.path[i].X == x && j.path[i].Y == y {
			if j.complet == false && j.previous != j.path[i] {
				fmt.Fprintf(os.Stderr, "Journey is complet\n")
				j.complet = true
			}
			return j.path[i]
		}
	}
	p := &Point{x, y}
	fmt.Fprintf(os.Stderr, "Append target to %d-%d\n", x, y)
	j.path = append(j.path, p)
	return p
}

func (j *Journey) SetPreviousTarget(t *Point) {
	j.previous = t
}

func (j *Journey) IsComplet() bool {
	//TODO
	return false
}

func main() {
	journey := Journey{nil, nil, false}
	for {
		var targetX, targetY int
		var speed string
		// nextCheckpointX: x position of the next check point
		// nextCheckpointY: y position of the next check point
		// nextCheckpointDist: distance to the next checkpoint
		// nextCheckpointAngle: angle between your pod orientation and the direction of the next checkpoint
		var x, y, nextCheckpointX, nextCheckpointY, nextCheckpointDist, nextCheckpointAngle int
		fmt.Scan(&x, &y, &nextCheckpointX, &nextCheckpointY, &nextCheckpointDist, &nextCheckpointAngle)
		currentTarget := journey.GetCurrentTarget(nextCheckpointX, nextCheckpointY)
		fmt.Fprintf(os.Stderr, "Current target is :%d-%d\n", currentTarget.X, currentTarget.Y)
		var opponentX, opponentY int
		fmt.Scan(&opponentX, &opponentY)

		//TODO Compute nearest target
		targetX = nextCheckpointX
		targetY = nextCheckpointY

		if journey.IsComplet() {
			//TODO improve
			speed = "100"
		} else {
			//TODO improve
			speed = "100"
			if nextCheckpointDist < 1000 {
				speed = "30"
			}

			if nextCheckpointAngle > 90 || nextCheckpointAngle < -90 {
				speed = "5"
			}
		}
		fmt.Printf("%d %d %s\n", targetX, targetY, speed)

		journey.SetPreviousTarget(currentTarget)
	}
}
