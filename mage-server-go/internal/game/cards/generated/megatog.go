package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/effects"
)

func init() {
	cards.Register("Megatog", NewMegatog)
}

// NewMegatog creates a Megatog
// {4}{R}{R} - CREATURE
func NewMegatog(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Megatog")
	card.ManaCost = "{4}{R}{R}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"ATOG"}
	card.Power = "3"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddEffect(abilities.NewBoostEffect(3, 3)).
		AddEffect(abilities.NewGrantAbilityEffect("TrampleAbility", effects.DurationEndOfTurn)).
		Build()
	card.AddAbility(ability0)
	return card, nil
}
