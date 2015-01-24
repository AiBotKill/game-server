package main

import (
	"log"
	"time"
)

type bullet struct {
	entity
	Damage  float64
	Shooter *player
}

func NewBullet(location [2]float64, velocity [2]float64, damage float64, game *game, shooter *player) *bullet {
	b := &bullet{entity: *NewEntity(location, 2, game)}
	b.Velocity = velocity
	b.Damage = damage
	b.Shooter = shooter
	return b
}

func (b *bullet) update(dt time.Duration) {
	// Store old location before calling entity.update
	oldLocation := b.entity.Location
	b.entity.update(dt)

	// As game for all collision between old and new location
	var collisionVector [4]float64
	collisionVector[0] = oldLocation[0]
	collisionVector[1] = oldLocation[1]
	collisionVector[2] = b.entity.Location[0]
	collisionVector[3] = b.entity.Location[1]
	collisions := b.entity.Game.collision(collisionVector)

	// Damage all entities unfortunate enough to get in the way
	for _, e := range collisions {
		log.Println("bullet collisions", e)
	}
	log.Println("bullet updated:", collisionVector)
}
