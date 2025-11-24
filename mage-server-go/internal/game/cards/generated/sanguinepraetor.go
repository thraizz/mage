package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Sanguine Praetor", NewSanguinePraetor)
}

// NewSanguinePraetor creates a Sanguine Praetor
// {6}{B}{B} - CREATURE
func NewSanguinePraetor(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Sanguine Praetor")
	card.ManaCost = "{6}{B}{B}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"AVATAR", "PRAETOR"}
	card.Power = "7"
	card.Toughness = "5"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}