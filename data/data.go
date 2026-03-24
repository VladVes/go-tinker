package data

import (
	"github.com/VladVes/go-tinker/v2/data/entities"
	"github.com/VladVes/go-tinker/v2/data/schemas"
)

var (
	PostLikes = map[string]int64{}
	Logs      []schemas.LogEntry
	Orders    = make(map[string]entities.Order)
	Employee  = make(map[string]entities.Employee)
	Users     = make(map[string]entities.User)
)
