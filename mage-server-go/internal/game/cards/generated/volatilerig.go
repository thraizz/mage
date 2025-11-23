package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Volatile Rig", NewVolatileRig)
}

// NewVolatileRig creates a Volatile Rig
// {4} - ARTIFACT CREATURE
// Trample
func NewVolatileRig(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Volatile Rig")
	card.ManaCost = "{4}"
	card.Types = []string{"ARTIFACT", "CREATURE"}
	card.Subtypes = []string{"CONSTRUCT"}
	card.Power = "4"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordTrample)
	card.AddAbility(ability0)
	return card, nil
}
