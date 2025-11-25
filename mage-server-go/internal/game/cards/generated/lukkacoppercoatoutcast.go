package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Lukka Coppercoat Outcast", NewLukkaCoppercoatOutcast)
}

// NewLukkaCoppercoatOutcast creates a Lukka Coppercoat Outcast
// {3}{R}{R} - PLANESWALKER
func NewLukkaCoppercoatOutcast(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Lukka Coppercoat Outcast")
	card.ManaCost = "{3}{R}{R}"
	card.Types = []string{"PLANESWALKER"}
	card.Subtypes = []string{"LUKKA"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Loyalty = "5"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement triggered ability: LoyaltyAbility
	//   - Effect: LukkaCoppercoatOutcastPolymorphEffect()
	// card.AddAbility(ability0)
	return card, nil
}
