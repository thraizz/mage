package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Scourge Of Kher Ridges", NewScourgeOfKherRidges)
}

// NewScourgeOfKherRidges creates a Scourge Of Kher Ridges
// {6}{R}{R} - CREATURE
// Flying
func NewScourgeOfKherRidges(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Scourge Of Kher Ridges")
	card.ManaCost = "{6}{R}{R}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"DRAGON"}
	card.Power = "6"
	card.Toughness = "6"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	// TODO: Implement activated ability with unmapped effects
	//   - DamageAllEffect(2, filter)
	// card.AddAbility(ability1)
	// TODO: Implement activated ability with unmapped effects
	//   - DamageAllEffect(6, filter2)
	// card.AddAbility(ability2)
	return card, nil
}
