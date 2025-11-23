package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/counters"
	"github.com/magefree/mage-server-go/internal/game/effects"
)

func init() {
	cards.Register("Azlask The Swelling Scourge", NewAzlaskTheSwellingScourge)
}

// NewAzlaskTheSwellingScourge creates a Azlask The Swelling Scourge
// {3} - CREATURE
func NewAzlaskTheSwellingScourge(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Azlask The Swelling Scourge")
	card.ManaCost = "{3}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"ELDRAZI"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "2"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddEffect(abilities.NewBoostEffect(SourceControllerCountersCount.EXPERIENCE, SourceControllerCountersCount.EXPERIENCE)).
		AddEffect(abilities.NewGrantAbilityEffect("IndestructibleAbility", effects.DurationEndOfTurn)).
		AddEffect(abilities.NewGrantAbilityEffect(new AnnihilatorAbility(1), filter2)).
		Build()
	card.AddAbility(ability0)
	return card, nil
}