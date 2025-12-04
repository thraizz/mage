package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Gossamer Chains", NewGossamerChains)
}

// NewGossamerChains creates a Gossamer Chains
// {W}{W} - ENCHANTMENT
func NewGossamerChains(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Gossamer Chains")
	card.ManaCost = "{W}{W}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement activated ability with unmapped effects
	//   - PreventDamageByTargetEffect()
	// card.AddAbility(ability0)
	return card, nil
}
