package jwt

import (
	"sync"
	"time"
)

type TokenBlacklist struct {
	mu *sync.Mutex
	tkn map[string]time.Time
}

func NewTokenBlacklist() *TokenBlacklist {
	return &TokenBlacklist{
		mu: &sync.Mutex{},
		tkn: make(map[string]time.Time),
	}
}

func (bl *TokenBlacklist) AddTokenBlacklist(token string, expiresAt time.Time) {
	bl.mu.Lock()
	bl.tkn[token] = expiresAt
	bl.mu.Unlock()
}

func (bl *TokenBlacklist) IsBlacklisted(token string) bool {
	bl.mu.Lock()
	defer bl.mu.Unlock()
	_, exist := bl.tkn[token]
	return  exist
}