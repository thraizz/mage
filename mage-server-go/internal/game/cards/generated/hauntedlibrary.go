package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/token"
)

func init() {
	cards.Register("Haunted Library", NewHauntedLibrary)
}

// NewHauntedLibrary creates a Haunted Library
// {1}{W} - ENCHANTMENT
func NewHauntedLibrary(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Haunted Library")
	card.ManaCost = "{1}{W}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - DoIfCostPaid(                 new CreateTokenEffect(new SpiritW...)
	// card.AddAbility(ability0)
	return card, nil
}
