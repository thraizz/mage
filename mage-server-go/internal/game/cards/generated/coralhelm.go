package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Coral Helm", NewCoralHelm)
}

// NewCoralHelm creates a Coral Helm
// {3} - ARTIFACT
func NewCoralHelm(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Coral Helm")
	card.ManaCost = "{3}"
	card.Types = []string{"ARTIFACT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddManaCost("{3}").
		AddEffect(abilities.NewBoostEffect(2, 2)).
		Build()
	card.AddAbility(ability0)
	return card, nil
}
