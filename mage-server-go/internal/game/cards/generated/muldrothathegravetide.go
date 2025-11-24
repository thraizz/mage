package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Muldrotha The Gravetide", NewMuldrothaTheGravetide)
}

// NewMuldrothaTheGravetide creates a Muldrotha The Gravetide
// {3}{B}{G}{U} - CREATURE
func NewMuldrothaTheGravetide(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Muldrotha The Gravetide")
	card.ManaCost = "{3}{B}{G}{U}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"ELEMENTAL", "AVATAR"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "6"
	card.Toughness = "6"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}