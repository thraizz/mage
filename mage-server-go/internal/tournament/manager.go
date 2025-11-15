package tournament

import (
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// TournamentState represents the state of a tournament
type TournamentState int

const (
	TournamentStateWaiting TournamentState = iota
	TournamentStateStarting
	TournamentStateDrafting
	TournamentStateConstructing
	TournamentStateInProgress
	TournamentStateFinished
)

func (s TournamentState) String() string {
	switch s {
	case TournamentStateWaiting:
		return "WAITING"
	case TournamentStateStarting:
		return "STARTING"
	case TournamentStateDrafting:
		return "DRAFTING"
	case TournamentStateConstructing:
		return "CONSTRUCTING"
	case TournamentStateInProgress:
		return "IN_PROGRESS"
	case TournamentStateFinished:
		return "FINISHED"
	default:
		return "UNKNOWN"
	}
}

// TournamentType represents the type of tournament
type TournamentType string

const (
	TournamentTypeSwiss       TournamentType = "SWISS"
	TournamentTypeElimination TournamentType = "ELIMINATION"
)

// Player represents a tournament participant
type Player struct {
	Name       string
	Points     int
	Wins       int
	Losses     int
	Draws      int
	Eliminated bool
	Quit       bool
}

// Pairing represents a match pairing in a round
type Pairing struct {
	Player1     string
	Player2     string
	TableID     string
	Winner      string
	Player1Wins int
	Player2Wins int
	Draws       int
}

// Round represents a tournament round
type Round struct {
	Number   int
	Pairings []*Pairing
	Started  bool
	Finished bool
}

// PlayerSnapshot captures tournament player data for external use.
type PlayerSnapshot struct {
	Name       string
	Points     int
	Wins       int
	Losses     int
	Draws      int
	Eliminated bool
	Quit       bool
}

// PairingSnapshot captures pairing data for external use.
type PairingSnapshot struct {
	Player1     string
	Player2     string
	TableID     string
	Winner      string
	Player1Wins int
	Player2Wins int
	Draws       int
}

// RoundSnapshot captures round data for external use.
type RoundSnapshot struct {
	Number   int
	Started  bool
	Finished bool
	Pairings []PairingSnapshot
}

// TournamentSnapshot captures a consistent view of a tournament.
type TournamentSnapshot struct {
	ID             string
	Name           string
	Type           string
	State          TournamentState
	ControllerName string
	RoomID         string
	Players        []PlayerSnapshot
	Rounds         []RoundSnapshot
	CurrentRound   int
	NumRounds      int
	WinsRequired   int
	CreateTime     time.Time
	StartTime      *time.Time
	EndTime        *time.Time
}

// Tournament represents a tournament
type Tournament struct {
	ID                string
	Name              string
	Type              TournamentType
	TournamentTypeStr string // "Constructed", "Booster Draft", "Sealed"
	State             TournamentState
	ControllerName    string
	RoomID            string
	Players           map[string]*Player
	PlayerOrder       []string // Maintains insertion order
	Rounds            []*Round
	CurrentRound      int
	NumRounds         int
	WinsRequired      int
	CreateTime        time.Time
	StartTime         *time.Time
	EndTime           *time.Time
	Winner            string
	Watchers          map[string]bool
	mu                sync.RWMutex
}

// NewTournament creates a new tournament
func NewTournament(name, tournamentType, controllerName, roomID string, numRounds, winsRequired int) *Tournament {
	return &Tournament{
		ID:                uuid.New().String(),
		Name:              name,
		Type:              TournamentTypeSwiss,
		TournamentTypeStr: tournamentType,
		State:             TournamentStateWaiting,
		ControllerName:    controllerName,
		RoomID:            roomID,
		Players:           make(map[string]*Player),
		PlayerOrder:       make([]string, 0),
		Rounds:            make([]*Round, 0),
		CurrentRound:      0,
		NumRounds:         numRounds,
		WinsRequired:      winsRequired,
		CreateTime:        time.Now(),
		Watchers:          make(map[string]bool),
	}
}

// AddPlayer adds a player to the tournament
func (t *Tournament) AddPlayer(playerName string) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.State != TournamentStateWaiting {
		return fmt.Errorf("tournament already started")
	}

	if _, exists := t.Players[playerName]; exists {
		return fmt.Errorf("player already joined")
	}

	t.Players[playerName] = &Player{
		Name:   playerName,
		Points: 0,
		Wins:   0,
		Losses: 0,
		Draws:  0,
	}
	t.PlayerOrder = append(t.PlayerOrder, playerName)

	return nil
}

// RemovePlayer removes a player from the tournament
func (t *Tournament) RemovePlayer(playerName string) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.State != TournamentStateWaiting {
		return fmt.Errorf("tournament already started")
	}

	if _, exists := t.Players[playerName]; !exists {
		return fmt.Errorf("player not found")
	}

	delete(t.Players, playerName)

	// Remove from order
	for i, name := range t.PlayerOrder {
		if name == playerName {
			t.PlayerOrder = append(t.PlayerOrder[:i], t.PlayerOrder[i+1:]...)
			break
		}
	}

	return nil
}

// GetPlayerCount returns the number of players
func (t *Tournament) GetPlayerCount() int {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return len(t.Players)
}

// AddWatcher registers a watcher for the tournament.
func (t *Tournament) AddWatcher(playerName string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.Watchers[playerName] = true
}

// RemoveWatcher removes a watcher from the tournament.
func (t *Tournament) RemoveWatcher(playerName string) bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	if _, exists := t.Watchers[playerName]; exists {
		delete(t.Watchers, playerName)
		return true
	}
	return false
}

// GetWatchers returns all watchers currently observing the tournament.
func (t *Tournament) GetWatchers() []string {
	t.mu.RLock()
	defer t.mu.RUnlock()

	watchers := make([]string, 0, len(t.Watchers))
	for watcher := range t.Watchers {
		watchers = append(watchers, watcher)
	}
	return watchers
}

// GetPlayers returns all players
func (t *Tournament) GetPlayers() []*Player {
	t.mu.RLock()
	defer t.mu.RUnlock()

	players := make([]*Player, 0, len(t.Players))
	for _, name := range t.PlayerOrder {
		if player, ok := t.Players[name]; ok {
			players = append(players, player)
		}
	}
	return players
}

// QuitPlayer marks a player as having quit an active tournament.
func (t *Tournament) QuitPlayer(playerName string) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	player, exists := t.Players[playerName]
	if !exists {
		return fmt.Errorf("player not found")
	}

	player.Quit = true
	return nil
}

// SetState sets the tournament state
func (t *Tournament) SetState(state TournamentState) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.State = state

	if state == TournamentStateInProgress && t.StartTime == nil {
		now := time.Now()
		t.StartTime = &now
	} else if state == TournamentStateFinished {
		now := time.Now()
		t.EndTime = &now
	}
}

// GetState returns the current tournament state
func (t *Tournament) GetState() TournamentState {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.State
}

// CreateRound creates a new round with pairings
func (t *Tournament) CreateRound() *Round {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.CurrentRound++
	round := &Round{
		Number:   t.CurrentRound,
		Pairings: t.generatePairings(),
		Started:  false,
		Finished: false,
	}

	t.Rounds = append(t.Rounds, round)
	return round
}

// generatePairings generates pairings for the current round (Swiss pairing algorithm)
func (t *Tournament) generatePairings() []*Pairing {
	// Simple Swiss pairing: pair players with similar points
	activePlayers := make([]*Player, 0)
	for _, player := range t.Players {
		if !player.Eliminated && !player.Quit {
			activePlayers = append(activePlayers, player)
		}
	}

	// Sort by points (simplified - real implementation would be more complex)
	// For now, just pair sequentially
	pairings := make([]*Pairing, 0)
	for i := 0; i < len(activePlayers)-1; i += 2 {
		pairing := &Pairing{
			Player1: activePlayers[i].Name,
			Player2: activePlayers[i+1].Name,
		}
		pairings = append(pairings, pairing)
	}

	// Handle bye if odd number of players
	if len(activePlayers)%2 == 1 {
		// Last player gets a bye (automatic win)
		lastPlayer := activePlayers[len(activePlayers)-1]
		lastPlayer.Points += 3
		lastPlayer.Wins++
	}

	return pairings
}

// RecordMatchResult records the result of a match
func (t *Tournament) RecordMatchResult(roundNum int, player1, player2, winner string, player1Wins, player2Wins int) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if roundNum <= 0 || roundNum > len(t.Rounds) {
		return fmt.Errorf("invalid round number")
	}

	round := t.Rounds[roundNum-1]

	// Find the pairing
	for _, pairing := range round.Pairings {
		if (pairing.Player1 == player1 && pairing.Player2 == player2) ||
			(pairing.Player1 == player2 && pairing.Player2 == player1) {
			pairing.Winner = winner
			pairing.Player1Wins = player1Wins
			pairing.Player2Wins = player2Wins

			// Update player stats
			if winner == player1 {
				t.Players[player1].Wins++
				t.Players[player1].Points += 3
				t.Players[player2].Losses++
			} else if winner == player2 {
				t.Players[player2].Wins++
				t.Players[player2].Points += 3
				t.Players[player1].Losses++
			} else {
				// Draw
				t.Players[player1].Draws++
				t.Players[player1].Points += 1
				t.Players[player2].Draws++
				t.Players[player2].Points += 1
			}

			return nil
		}
	}

	return fmt.Errorf("pairing not found")
}

// IsController checks if the given player is the tournament controller
func (t *Tournament) IsController(playerName string) bool {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.ControllerName == playerName
}

// Start transitions the tournament into progress and creates the first round.
func (t *Tournament) Start() error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.State != TournamentStateWaiting {
		return fmt.Errorf("tournament already started")
	}

	if len(t.Players) < 2 {
		return fmt.Errorf("not enough players")
	}

	now := time.Now()
	if t.StartTime == nil {
		t.StartTime = &now
	}

	t.State = TournamentStateInProgress
	t.CurrentRound = 0

	t.CurrentRound++
	round := &Round{
		Number:   t.CurrentRound,
		Pairings: t.generatePairings(),
		Started:  true,
		Finished: false,
	}
	t.Rounds = append(t.Rounds, round)

	return nil
}

// Snapshot returns a consistent copy of the tournament state.
func (t *Tournament) Snapshot() TournamentSnapshot {
	t.mu.RLock()
	defer t.mu.RUnlock()

	players := make([]PlayerSnapshot, 0, len(t.PlayerOrder))
	for _, name := range t.PlayerOrder {
		if player, ok := t.Players[name]; ok {
			players = append(players, PlayerSnapshot{
				Name:       player.Name,
				Points:     player.Points,
				Wins:       player.Wins,
				Losses:     player.Losses,
				Draws:      player.Draws,
				Eliminated: player.Eliminated,
				Quit:       player.Quit,
			})
		}
	}

	rounds := make([]RoundSnapshot, 0, len(t.Rounds))
	for _, r := range t.Rounds {
		pairings := make([]PairingSnapshot, 0, len(r.Pairings))
		for _, p := range r.Pairings {
			pairings = append(pairings, PairingSnapshot{
				Player1:     p.Player1,
				Player2:     p.Player2,
				TableID:     p.TableID,
				Winner:      p.Winner,
				Player1Wins: p.Player1Wins,
				Player2Wins: p.Player2Wins,
				Draws:       p.Draws,
			})
		}

		rounds = append(rounds, RoundSnapshot{
			Number:   r.Number,
			Started:  r.Started,
			Finished: r.Finished,
			Pairings: pairings,
		})
	}

	return TournamentSnapshot{
		ID:             t.ID,
		Name:           t.Name,
		Type:           t.TournamentTypeStr,
		State:          t.State,
		ControllerName: t.ControllerName,
		RoomID:         t.RoomID,
		Players:        players,
		Rounds:         rounds,
		CurrentRound:   t.CurrentRound,
		NumRounds:      t.NumRounds,
		WinsRequired:   t.WinsRequired,
		CreateTime:     t.CreateTime,
		StartTime:      cloneTime(t.StartTime),
		EndTime:        cloneTime(t.EndTime),
	}
}

func cloneTime(src *time.Time) *time.Time {
	if src == nil {
		return nil
	}
	cp := *src
	return &cp
}

// Manager manages tournaments
type Manager struct {
	tournaments map[string]*Tournament
	mu          sync.RWMutex
	logger      *zap.Logger
}

// NewManager creates a new tournament manager
func NewManager(logger *zap.Logger) *Manager {
	return &Manager{
		tournaments: make(map[string]*Tournament),
		logger:      logger,
	}
}

// CreateTournament creates a new tournament
func (m *Manager) CreateTournament(name, tournamentType, controllerName, roomID string, numRounds, winsRequired int) *Tournament {
	m.mu.Lock()
	defer m.mu.Unlock()

	tournament := NewTournament(name, tournamentType, controllerName, roomID, numRounds, winsRequired)
	m.tournaments[tournament.ID] = tournament

	m.logger.Info("tournament created",
		zap.String("tournament_id", tournament.ID),
		zap.String("name", name),
		zap.String("type", tournamentType),
		zap.String("controller", controllerName),
	)

	return tournament
}

// GetTournament retrieves a tournament by ID
func (m *Manager) GetTournament(tournamentID string) (*Tournament, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	tournament, ok := m.tournaments[tournamentID]
	return tournament, ok
}

// RemoveTournament removes a tournament
func (m *Manager) RemoveTournament(tournamentID string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.tournaments, tournamentID)

	m.logger.Info("tournament removed", zap.String("tournament_id", tournamentID))
}

// GetAllTournaments returns all tournaments
func (m *Manager) GetAllTournaments() []*Tournament {
	m.mu.RLock()
	defer m.mu.RUnlock()

	tournaments := make([]*Tournament, 0, len(m.tournaments))
	for _, tournament := range m.tournaments {
		tournaments = append(tournaments, tournament)
	}
	return tournaments
}

// GetActiveTournamentCount returns the count of active tournaments
func (m *Manager) GetActiveTournamentCount() int {
	m.mu.RLock()
	defer m.mu.RUnlock()

	count := 0
	for _, tournament := range m.tournaments {
		if tournament.State != TournamentStateFinished {
			count++
		}
	}
	return count
}
