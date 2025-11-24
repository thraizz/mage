package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Grixis Charm", NewGrixisCharm)
}

// NewGrixisCharm creates a Grixis Charm
// {U}{B}{R} - INSTANT
func NewGrixisCharm(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Grixis Charm")
	card.ManaCost = "{U}{B}{R}"
	card.Types = []string{"INSTANT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewBoostEffect(-4, -4)).
		AddEffect(abilities.NewBoostEffect(2, 0, false)).
		// TODO: ReturnToHandTargetEffect with complex parameters
		AddEffect(abilities.NewBoostEffect(-4, -4)).
		AddEffect(abilities.NewBoostEffect(2, 0, false)).
		AddTarget(abilities.NewPermanentTargetFilter()).
		AddTarget(abilities.NewCreatureTargetFilter()).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
