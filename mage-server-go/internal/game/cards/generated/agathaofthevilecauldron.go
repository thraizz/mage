package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/effects"
)

func init() {
	cards.Register("Agatha Of The Vile Cauldron", NewAgathaOfTheVileCauldron)
}

// NewAgathaOfTheVileCauldron creates a Agatha Of The Vile Cauldron
// {R}{G} - CREATURE
func NewAgathaOfTheVileCauldron(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Agatha Of The Vile Cauldron")
	card.ManaCost = "{R}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"HUMAN", "WARLOCK"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "1"
	card.Toughness = "1"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddEffect(abilities.NewBoostEffect(1, 1, true)).
		AddEffect(abilities.NewGrantAbilityEffect("TrampleAbility", effects.DurationEndOfTurn)).
		AddEffect(abilities.NewGrantAbilityEffect("HasteAbility", effects.DurationEndOfTurn)).
		Build()
	card.AddAbility(ability0)
	return card, nil
}
