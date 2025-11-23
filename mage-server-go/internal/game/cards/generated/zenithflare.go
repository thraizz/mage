package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Zenith Flare", NewZenithFlare)
}

// NewZenithFlare creates a Zenith Flare
// {2}{R}{W} - INSTANT
func NewZenithFlare(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Zenith Flare")
	card.ManaCost = "{2}{R}{W}"
	card.Types = []string{"INSTANT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewDamageEffect(xValue)).
		AddEffect(abilities.NewGainLifeEffect(xValue)).
		AddTarget(abilities.NewAnyTargetFilter()).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
