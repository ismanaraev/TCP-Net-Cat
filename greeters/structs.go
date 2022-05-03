package chat

import "sync"

type User struct {
	Username string
}

var (
	Users   int
	Usernum sync.Mutex
)
