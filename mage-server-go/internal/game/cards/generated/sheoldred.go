package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Sheoldred", NewSheoldred)
}

// NewSheoldred creates a Sheoldred
// {3}{B}{B} - CREATURE
func NewSheoldred(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Sheoldred")
	card.ManaCost = "{3}{B}{B}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"PHYREXIAN", "PRAETOR"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "4"
	card.Toughness = "5"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
