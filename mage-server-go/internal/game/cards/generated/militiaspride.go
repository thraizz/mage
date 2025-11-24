package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/token"
)

func init() {
	cards.Register("Militias Pride", NewMilitiasPride)
}

// NewMilitiasPride creates a Militias Pride
// {1}{W} - KINDRED ENCHANTMENT
func NewMilitiasPride(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Militias Pride")
	card.ManaCost = "{1}{W}"
	card.Types = []string{"KINDRED", "ENCHANTMENT"}
	card.Subtypes = []string{"KITHKIN"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - DoIfCostPaid(new CreateTokenEffect(new KithkinSoldierToken(), 1...)
	// card.AddAbility(ability0)
	return card, nil
}