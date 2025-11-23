package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Rootwater Mystic", NewRootwaterMystic)
}

// NewRootwaterMystic creates a Rootwater Mystic
// {U} - CREATURE
func NewRootwaterMystic(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Rootwater Mystic")
	card.ManaCost = "{U}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"MERFOLK", "WIZARD"}
	card.Power = "1"
	card.Toughness = "1"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
