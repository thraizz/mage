package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Spined Thopter", NewSpinedThopter)
}

// NewSpinedThopter creates a Spined Thopter
// {2}{U/P} - ARTIFACT CREATURE
// Flying
func NewSpinedThopter(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Spined Thopter")
	card.ManaCost = "{2}{U/P}"
	card.Types = []string{"ARTIFACT", "CREATURE"}
	card.Subtypes = []string{"PHYREXIAN", "THOPTER"}
	card.Power = "2"
	card.Toughness = "1"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	return card, nil
}
