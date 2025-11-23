package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Synth Eradicator", NewSynthEradicator)
}

// NewSynthEradicator creates a Synth Eradicator
// {2}{R} - ARTIFACT CREATURE
// Haste
func NewSynthEradicator(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Synth Eradicator")
	card.ManaCost = "{2}{R}"
	card.Types = []string{"ARTIFACT", "CREATURE"}
	card.Subtypes = []string{"SYNTH", "SOLDIER"}
	card.Power = "3"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordHaste)
	card.AddAbility(ability0)
	ability1 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddTapCost().
		AddEffect(abilities.NewDamageEffect(3)).
		Build()
	card.AddAbility(ability1)
	return card, nil
}
