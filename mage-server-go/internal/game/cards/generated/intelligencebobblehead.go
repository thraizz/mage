package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Intelligence Bobblehead", NewIntelligenceBobblehead)
}

// NewIntelligenceBobblehead creates a Intelligence Bobblehead
// {3} - ARTIFACT
func NewIntelligenceBobblehead(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Intelligence Bobblehead")
	card.ManaCost = "{3}"
	card.Types = []string{"ARTIFACT"}
	card.Subtypes = []string{"BOBBLEHEAD"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddManaCost("{5}").
		AddTapCost().
		AddEffect(abilities.NewDrawCardsEffect(xValue)).
		Build()
	card.AddAbility(ability0)
	return card, nil
}
