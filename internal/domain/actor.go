package domain

import "time"

type Actor struct {
	Id        int
	Name      string
	Gender    int // http://en.wikipedia.org/wiki/ISO_5218
	BirthDate time.Time
}
