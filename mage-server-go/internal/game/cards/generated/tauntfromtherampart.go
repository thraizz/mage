package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Taunt From The Rampart", NewTauntFromTheRampart)
}

// NewTauntFromTheRampart creates a Taunt From The Rampart
// {3}{R}{W} - SORCERY
func NewTauntFromTheRampart(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Taunt From The Rampart")
	card.ManaCost = "{3}{R}{W}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
