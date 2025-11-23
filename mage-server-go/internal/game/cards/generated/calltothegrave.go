package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Call To The Grave", NewCallToTheGrave)
}

// NewCallToTheGrave creates a Call To The Grave
// {4}{B} - ENCHANTMENT
func NewCallToTheGrave(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Call To The Grave")
	card.ManaCost = "{4}{B}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - SacrificeSourceEffect()
	// card.AddAbility(ability0)
	return card, nil
}
