package game

import (
	"crypto/rand"
	"encoding/hex"
	"sync"
)

type Game struct {
	mutex sync.Mutex
	state Snapshot // guarded by mutex, except Id which never changes
	subs  map[chan Snapshot]struct{}
}

func New(id string, prompt string) *Game {
	return &Game{
		state: Snapshot{Id: id, Prompt: prompt, Players: []Player{}},
		subs:  map[chan Snapshot]struct{}{},
	}
}

func (g *Game) Id() string { return g.state.Id }

// A Snapshot is an immutable view of the Game state.
type Snapshot struct {
	Id      string   `json:"id"`
	Players []Player `json:"players"`
	Version uint64   `json:"version"`
	Prompt  string   `json:"prompt"`
}

type Player struct {
	Name  string `json:"name"`
	Token string `json:"-"`
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
	g.mutex.Lock()
	defer g.mutex.Unlock()
	p := Player{Name: name, Token: generateSecureToken(16)}
	g.state.Players = append(g.state.Players, p)
	g.publish()
	return p
}

func (g *Game) PlayerByToken(token string) (Player, bool) {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	for _, p := range g.state.Players {
		if p.Token == token {
			return p, true
		}
	}
	return Player{}, false
}

// publish must be called with g.mutex held.
func (g *Game) publish() {
	g.state.Version += 1
	s := g.state
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
	return g.state
}

// Subscribe returns a channel that receives the current state immediately and
// again after every change. Call the returned function to unsubscribe, which
// also closes the channel.
func (g *Game) Subscribe() (<-chan Snapshot, func()) {
	ch := make(chan Snapshot, 1)
	g.mutex.Lock()
	defer g.mutex.Unlock()
	ch <- g.state
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
