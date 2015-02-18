package main

import "log"

func SortCollisions(collisions []*collision, position *Vector) []*collision {
	// Sort collisions
	for i := 0; i < len(collisions)-1; i++ {
		for j := i + 1; j < len(collisions); j++ {
			v1 := collisions[i].Position.Sub(position)
			v2 := collisions[j].Position.Sub(position)
			log.Println(v1, v2, v1.Length(), v2.Length())
			if v1.Length() > v2.Length() {
				collisions[i], collisions[j] = collisions[j], collisions[i]
			}
		}
	}
	return collisions
}
