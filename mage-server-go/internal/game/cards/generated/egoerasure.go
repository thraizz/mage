package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Ego Erasure", NewEgoErasure)
}

// NewEgoErasure creates a Ego Erasure
// {2}{U} - KINDRED INSTANT
func NewEgoErasure(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Ego Erasure")
	card.ManaCost = "{2}{U}"
	card.Types = []string{"KINDRED", "INSTANT"}
	card.Subtypes = []string{"SHAPESHIFTER"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddTarget(abilities.NewTargetRequirement(1, 1, abilities.NewPlayerTargetFilter())).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
