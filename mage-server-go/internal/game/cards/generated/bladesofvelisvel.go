package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Blades Of Velis Vel", NewBladesOfVelisVel)
}

// NewBladesOfVelisVel creates a Blades Of Velis Vel
// {1}{R} - KINDRED INSTANT
func NewBladesOfVelisVel(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Blades Of Velis Vel")
	card.ManaCost = "{1}{R}"
	card.Types = []string{"KINDRED", "INSTANT"}
	card.Subtypes = []string{"SHAPESHIFTER"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewBoostEffect(2, 0)).
		AddTarget(abilities.NewTargetRequirement(0, 2, abilities.NewCreatureTargetFilter())).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
