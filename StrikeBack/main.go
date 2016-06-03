package main

import "fmt"

//import "os"

/**
 * Auto-generated code below aims at helping you parse
 * the standard input according to the problem statement.
 **/

func main() {
	boost_used := false
	for {
		// nextCheckpointX: x position of the next check point
		// nextCheckpointY: y position of the next check point
		// nextCheckpointDist: distance to the next checkpoint
		// nextCheckpointAngle: angle between your pod orientation and the direction of the next checkpoint
		var x, y, nextCheckpointX, nextCheckpointY, nextCheckpointDist, nextCheckpointAngle int
		fmt.Scan(&x, &y, &nextCheckpointX, &nextCheckpointY, &nextCheckpointDist, &nextCheckpointAngle)

		var opponentX, opponentY int
		fmt.Scan(&opponentX, &opponentY)

		speed := "100"

		if nextCheckpointDist > 10000 && boost_used == false && nextCheckpointAngle < 5 || nextCheckpointAngle > -5 {
			speed = "BOOST"
		}

		if nextCheckpointDist < 1000 {
			speed = "30"
		}

		if nextCheckpointAngle > 90 || nextCheckpointAngle < -90 {
			speed = "5"
		}

		if speed == "BOOST" {
			boost_used = true
		}

		// fmt.Fprintln(os.Stderr, "Debug messages...")

		// You have to output the target position
		// followed by the power (0 <= thrust <= 100)
		// i.e.: "x y thrust"
		fmt.Printf("%d %d %s\n", nextCheckpointX, nextCheckpointY, speed)
	}
}
