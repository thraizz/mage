package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Magnivore", NewMagnivore)
}

// NewMagnivore creates a Magnivore
// {2}{R}{R} - CREATURE
// Haste
func NewMagnivore(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Magnivore")
	card.ManaCost = "{2}{R}{R}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"LHURGOYF"}
	card.Power = "0"
	card.Toughness = "0"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordHaste)
	card.AddAbility(ability0)
	return card, nil
}
