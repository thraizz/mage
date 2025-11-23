package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Steelshaper Apprentice", NewSteelshaperApprentice)
}

// NewSteelshaperApprentice creates a Steelshaper Apprentice
// {2}{W}{W} - CREATURE
func NewSteelshaperApprentice(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Steelshaper Apprentice")
	card.ManaCost = "{2}{W}{W}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"HUMAN", "SOLDIER"}
	card.Power = "1"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddTapCost().
		AddEffect(abilities.NewSearchLibraryPutInHandEffect(abilities.NewTargetRequirement(0, 1, abilities.NewAnyTargetFilter()), true)).
		Build()
	card.AddAbility(ability0)
	return card, nil
}
