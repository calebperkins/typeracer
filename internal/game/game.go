package game

import (
	"crypto/rand"
	"encoding/hex"
	"sync"
)

type Game struct {
	Id      string
	Prompt  string
	version uint64
	players []Player
	mutex   sync.Mutex
	subs    map[chan Snapshot]struct{}
}

func New(id string, prompt string) *Game {
	return &Game{Id: id, Prompt: prompt, players: []Player{}, subs: map[chan Snapshot]struct{}{}}
}

// A Snapshot is an immutable snapshot of the Game state.
type Snapshot struct {
	Id      string   `json:"id"`
	Players []Player `json:"players"`
	Version uint64   `json:"version"`
}

type Player struct {
	Name  string `json:"name"`
	Token string
}

func generateSecureToken(length int) string {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	return hex.EncodeToString(b)
}

// Advance the game state, for debugging only
func (g *Game) Refresh() {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	g.publish()
}

func (g *Game) Join(name string) Player {
	p := Player{Name: name, Token: generateSecureToken(16)}
	g.players = append(g.players, p)
	g.publish()
	return p
}

func (g *Game) PlayerByToken(token string) *Player {
	for _, p := range g.players {
		if p.Token == token {
			return &p
		}
	}
	return nil
}

// publish must be called with g.mutex held.
func (g *Game) publish() {
	g.version += 1
	s := g.snapshot()
	for ch := range g.subs {
		select {
		case ch <- s:
		default:
			// drop the stale snapshot so we can send a fresh one
			select {
			case <-ch:
			default:
			}
			ch <- s // cannot block: only publish sends, and it holds g.mutex
		}
	}
}

func (g *Game) Snapshot() Snapshot {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	return g.snapshot()
}

// snapshot must be called with g.mutex held.
func (g *Game) snapshot() Snapshot {
	return Snapshot{Id: g.Id, Players: g.players, Version: g.version}
}

// Subscribe returns a channel that receives the current state immediately and
// again after every change. Call the returned function to unsubscribe, which
// also closes the channel.
func (g *Game) Subscribe() (<-chan Snapshot, func()) {
	ch := make(chan Snapshot, 1)
	g.mutex.Lock()
	defer g.mutex.Unlock()
	ch <- g.snapshot()
	g.subs[ch] = struct{}{}
	return ch, func() {
		g.mutex.Lock()
		defer g.mutex.Unlock()
		if _, ok := g.subs[ch]; ok {
			delete(g.subs, ch)
			close(ch)
		}
	}
}
