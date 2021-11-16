package models

import "time"

type AccountID string

//Account represents a banking account
type Account struct {
	ID        int
	Name      string
	CPF       string
	Secret    string
	Balance   float64
	CreatedAt time.Time
}
