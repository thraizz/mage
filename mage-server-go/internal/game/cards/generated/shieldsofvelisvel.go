package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Shields Of Velis Vel", NewShieldsOfVelisVel)
}

// NewShieldsOfVelisVel creates a Shields Of Velis Vel
// {W} - KINDRED INSTANT
func NewShieldsOfVelisVel(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Shields Of Velis Vel")
	card.ManaCost = "{W}"
	card.Types = []string{"KINDRED", "INSTANT"}
	card.Subtypes = []string{"SHAPESHIFTER"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddTarget(abilities.NewPlayerTargetFilter()).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}