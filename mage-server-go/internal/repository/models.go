package repository

import (
	"encoding/json"
	"time"
)

// User represents a user account
type User struct {
	ID              int64
	Name            string
	Password        string // Argon2id hash
	Email           string
	Active          bool
	CreatedAt       time.Time
	UpdatedAt       time.Time
	LockEndTime     *time.Time
	ChatLockEndTime *time.Time
}

// IsLocked checks if the user account is locked
func (u *User) IsLocked() bool {
	if u.LockEndTime == nil {
		return false
	}
	return time.Now().Before(*u.LockEndTime)
}

// IsChatLocked checks if the user is chat-locked
func (u *User) IsChatLocked() bool {
	if u.ChatLockEndTime == nil {
		return false
	}
	return time.Now().Before(*u.ChatLockEndTime)
}

// Deck represents a user's deck collection entry
type Deck struct {
	ID          int64
	UserID      int64
	Name        string
	Format      string
	Description string
	MainDeck    []string // Card names
	Sideboard   []string // Card names
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// MainDeckJSON returns the main deck as JSON string
func (d *Deck) MainDeckJSON() (string, error) {
	data, err := json.Marshal(d.MainDeck)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// SideboardJSON returns the sideboard as JSON string
func (d *Deck) SideboardJSON() (string, error) {
	data, err := json.Marshal(d.Sideboard)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// SetMainDeckFromJSON parses JSON string into MainDeck slice
func (d *Deck) SetMainDeckFromJSON(jsonStr string) error {
	if jsonStr == "" {
		d.MainDeck = []string{}
		return nil
	}
	return json.Unmarshal([]byte(jsonStr), &d.MainDeck)
}

// SetSideboardFromJSON parses JSON string into Sideboard slice
func (d *Deck) SetSideboardFromJSON(jsonStr string) error {
	if jsonStr == "" {
		d.Sideboard = []string{}
		return nil
	}
	return json.Unmarshal([]byte(jsonStr), &d.Sideboard)
}

// CardCount returns the total number of cards in the deck (main + sideboard)
func (d *Deck) CardCount() int {
	return len(d.MainDeck) + len(d.Sideboard)
}

// MainDeckCount returns the number of cards in the main deck
func (d *Deck) MainDeckCount() int {
	return len(d.MainDeck)
}

// SideboardCount returns the number of cards in the sideboard
func (d *Deck) SideboardCount() int {
	return len(d.Sideboard)
}

// MatchPlayerInfo represents a player's information in a match
type MatchPlayerInfo struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	Deck     string `json:"deck,omitempty"` // Deck name or deck list summary
	Result   string `json:"result"`         // "win", "loss", "draw", "concede"
}

// MatchHistory represents a completed game/match record
type MatchHistory struct {
	ID              int64
	GameID          string
	TableID         string
	TournamentID    string
	Players         []MatchPlayerInfo // Player information
	GameType        string            // "TwoPlayerDuel", "Commander", etc.
	StartTime       time.Time
	EndTime         time.Time
	DurationSeconds int
	WinnerID        *int64 // Nullable - may be no winner (draw, timeout)
	WinnerName      string // Denormalized for efficiency
	MatchOptions    string // JSON string of match configuration
	ReplayData      string // Compressed game log (optional)
	CreatedAt       time.Time
}

// PlayersJSON returns the players slice as JSON string
func (m *MatchHistory) PlayersJSON() (string, error) {
	data, err := json.Marshal(m.Players)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// SetPlayersFromJSON parses JSON string into Players slice
func (m *MatchHistory) SetPlayersFromJSON(jsonStr string) error {
	if jsonStr == "" {
		m.Players = []MatchPlayerInfo{}
		return nil
	}
	return json.Unmarshal([]byte(jsonStr), &m.Players)
}

// Duration returns the match duration as a time.Duration
func (m *MatchHistory) Duration() time.Duration {
	return time.Duration(m.DurationSeconds) * time.Second
}
