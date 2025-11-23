package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Mist Cloaked Herald", NewMistCloakedHerald)
}

// NewMistCloakedHerald creates a Mist Cloaked Herald
// {U} - CREATURE
func NewMistCloakedHerald(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Mist Cloaked Herald")
	card.ManaCost = "{U}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"MERFOLK", "WARRIOR"}
	card.Power = "1"
	card.Toughness = "1"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
