package main

import (
	"fmt"
	"math"
	"os"
)

type Point struct {
	X int
	Y int
}

func (p *Point) MovePointTo(target *Point, distance int) {
	current_dx := float64(p.X - target.X)
	current_dy := float64(p.Y - target.Y)
	current_h := math.Hypot(current_dx, current_dy)
	target_h := float64(distance)
	if current_h < target_h {
		p.X = target.X
		p.Y = target.Y
		return
	}

	ratio := target_h / current_h
	dx := current_dx * ratio
	dy := current_dy * ratio

	p.X -= int(dx)
	p.Y -= int(dy)
}

type Pod struct {
	currentPosition  *Point
	previousPosition *Point
}

type (p *Pod) SetCurrentPosition(x int, y int){
	p.previousPosition = p.currentPosition
	p.currentPosotion = &Point{x,y}
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
	pod := Pod{nil, nil}
	for {
		var speed string
		// nextCheckpointX: x position of the next check point
		// nextCheckpointY: y position of the next check point
		// nextCheckpointDist: distance to the next checkpoint
		// nextCheckpointAngle: angle between your pod orientation and the direction of the next checkpoint
		var x, y, nextCheckpointX, nextCheckpointY, nextCheckpointDist, nextCheckpointAngle int
		fmt.Scan(&x, &y, &nextCheckpointX, &nextCheckpointY, &nextCheckpointDist, &nextCheckpointAngle)

		pod.SetCurrentPosition(x,y)
		fmt.Fprintf(os.Stderr, "Current position is :%d-%d with %d-%d speed\n", nextCheckpoint.X, nextCheckpoint.Y)

		nextCheckpoint := journey.GetCurrentTarget(nextCheckpointX, nextCheckpointY)
		fmt.Fprintf(os.Stderr, "Current checkpoint is :%d-%d\n", nextCheckpoint.X, nextCheckpoint.Y)

		var opponentX, opponentY int
		fmt.Scan(&opponentX, &opponentY)

		podPos := Point{x, y}

		target := *nextCheckpoint
		target.MovePointTo(&podPos, 500)

		fmt.Fprintf(os.Stderr, "Current checkpoint is :%d-%d\n", nextCheckpoint.X, nextCheckpoint.Y)
		fmt.Fprintf(os.Stderr, "Current target is :%d-%d\n", target.X, target.Y)

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
		fmt.Printf("%d %d %s\n", target.X, target.Y, speed)

		journey.SetPreviousTarget(nextCheckpoint)
	}
}
