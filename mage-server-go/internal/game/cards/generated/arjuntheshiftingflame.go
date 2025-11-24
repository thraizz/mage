package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Arjun The Shifting Flame", NewArjunTheShiftingFlame)
}

// NewArjunTheShiftingFlame creates a Arjun The Shifting Flame
// {4}{U}{R} - CREATURE
// Flying
func NewArjunTheShiftingFlame(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Arjun The Shifting Flame")
	card.ManaCost = "{4}{U}{R}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"SPHINX", "WIZARD"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "5"
	card.Toughness = "5"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	return card, nil
}
