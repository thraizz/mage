package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Survival Of The Fittest", NewSurvivalOfTheFittest)
}

// NewSurvivalOfTheFittest creates a Survival Of The Fittest
// {1}{G} - ENCHANTMENT
func NewSurvivalOfTheFittest(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Survival Of The Fittest")
	card.ManaCost = "{1}{G}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddEffect(abilities.NewSearchLibraryPutInHandEffect(abilities.NewTargetRequirement(0, 1, abilities.NewCreatureTargetFilter()), true)).
		Build()
	card.AddAbility(ability0)
	return card, nil
}