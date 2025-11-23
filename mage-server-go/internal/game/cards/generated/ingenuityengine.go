package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Ingenuity Engine", NewIngenuityEngine)
}

// NewIngenuityEngine creates a Ingenuity Engine
// {7} - ARTIFACT
func NewIngenuityEngine(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Ingenuity Engine")
	card.ManaCost = "{7}"
	card.Types = []string{"ARTIFACT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddManaCost("{1}").
		AddTapCost().
		AddEffect(abilities.NewReturnToHandTargetEffect()).
		Build()
	card.AddAbility(ability0)
	return card, nil
}
