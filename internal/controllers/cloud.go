package controllers

import (
	"sync"
)

type Cloud struct {
	Mutex sync.Mutex
	HVs   map[string]*HV
}
