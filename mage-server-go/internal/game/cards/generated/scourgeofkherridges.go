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
	ability1 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddEffect(abilities.NewDamageEffect(2, filter)).
		Build()
	card.AddAbility(ability1)
	ability2 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddEffect(abilities.NewDamageEffect(6, filter2)).
		Build()
	card.AddAbility(ability2)
	return card, nil
}
