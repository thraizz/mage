package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Bedeck Bedazzle", NewBedeckBedazzle)
}

// NewBedeckBedazzle creates a Bedeck Bedazzle
// {B/R}{B/R} - INSTANT
func NewBedeckBedazzle(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Bedeck Bedazzle")
	card.ManaCost = "{B/R}{B/R}"
	card.Types = []string{"INSTANT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		// TODO: DamageTargetEffect with complex parameters
		AddEffect(abilities.NewBoostEffect(3, -3)).
		AddTarget(abilities.NewCreatureTargetFilter()).
		AddTarget(abilities.NewLandTargetFilter()).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
