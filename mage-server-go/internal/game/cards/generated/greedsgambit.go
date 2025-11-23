package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Greeds Gambit", NewGreedsGambit)
}

// NewGreedsGambit creates a Greeds Gambit
// {3}{B} - ENCHANTMENT
func NewGreedsGambit(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Greeds Gambit")
	card.ManaCost = "{3}{B}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - DiscardControllerEffect(1)
	//   - DiscardControllerEffect(3)
	// card.AddAbility(ability0)
	return card, nil
}
