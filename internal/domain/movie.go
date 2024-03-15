package domain

import "time"

type Movie struct {
	Id          int
	Title       string
	Description string
	ReleaseDate time.Time
	Rating      float64
	Actors      []*Actor
}
