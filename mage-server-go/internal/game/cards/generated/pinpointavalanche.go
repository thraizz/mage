package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Pinpoint Avalanche", NewPinpointAvalanche)
}

// NewPinpointAvalanche creates a Pinpoint Avalanche
// {3}{R}{R} - INSTANT
func NewPinpointAvalanche(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Pinpoint Avalanche")
	card.ManaCost = "{3}{R}{R}"
	card.Types = []string{"INSTANT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewDamageEffect(4, false)).
		AddTarget(abilities.NewCreatureTargetFilter()).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
