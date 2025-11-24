package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Deathbellow War Cry", NewDeathbellowWarCry)
}

// NewDeathbellowWarCry creates a Deathbellow War Cry
// {5}{R}{R}{R} - SORCERY
func NewDeathbellowWarCry(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Deathbellow War Cry")
	card.ManaCost = "{5}{R}{R}{R}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}