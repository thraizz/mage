package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Crush Underfoot", NewCrushUnderfoot)
}

// NewCrushUnderfoot creates a Crush Underfoot
// {1}{R} - KINDRED INSTANT
func NewCrushUnderfoot(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Crush Underfoot")
	card.ManaCost = "{1}{R}"
	card.Types = []string{"KINDRED", "INSTANT"}
	card.Subtypes = []string{"GIANT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddTarget(abilities.NewTargetRequirement(1, 1, abilities.NewCreatureTargetFilter())).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
