package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Bazaar Of Wonders", NewBazaarOfWonders)
}

// NewBazaarOfWonders creates a Bazaar Of Wonders
// {3}{U}{U} - ENCHANTMENT
func NewBazaarOfWonders(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Bazaar Of Wonders")
	card.ManaCost = "{3}{U}{U}"
	card.Types = []string{"ENCHANTMENT"}
	card.Supertypes = []string{"WORLD"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - ExileGraveyardAllPlayersEffect()
	// card.AddAbility(ability0)
	return card, nil
}
