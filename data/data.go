package data

import (
	"github.com/VladVes/go-tinker/v2/data/entities"
	"github.com/VladVes/go-tinker/v2/data/schemas"
)

var PostLikes = map[string]int64{}
var Logs []schemas.LogEntry
var Orders = make(map[string]entities.Order)
