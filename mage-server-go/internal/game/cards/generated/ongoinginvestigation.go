package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Ongoing Investigation", NewOngoingInvestigation)
}

// NewOngoingInvestigation creates a Ongoing Investigation
// {1}{U} - ENCHANTMENT
func NewOngoingInvestigation(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Ongoing Investigation")
	card.ManaCost = "{1}{U}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement activated ability with unmapped effects
	//   - InvestigateEffect()
	// card.AddAbility(ability0)
	return card, nil
}
