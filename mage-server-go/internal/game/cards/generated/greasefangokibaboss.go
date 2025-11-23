package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/effects"
)

func init() {
	cards.Register("Greasefang Okiba Boss", NewGreasefangOkibaBoss)
}

// NewGreasefangOkibaBoss creates a Greasefang Okiba Boss
// {1}{W}{B} - CREATURE
func NewGreasefangOkibaBoss(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Greasefang Okiba Boss")
	card.ManaCost = "{1}{W}{B}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"RAT", "PILOT"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "4"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewGrantAbilityEffect("HasteAbility", effects.DurationCustom)).
		AddEffect(abilities.NewReturnToHandTargetEffect()).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
