package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Warlock Class", NewWarlockClass)
}

// NewWarlockClass creates a Warlock Class
// {B} - ENCHANTMENT
func NewWarlockClass(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Warlock Class")
	card.ManaCost = "{B}"
	card.Types = []string{"ENCHANTMENT"}
	card.Subtypes = []string{"CLASS"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - LookLibraryAndPickControllerEffect(3, 1, PutCards.HAND, PutCards.GRAVEYARD)
	// card.AddAbility(ability0)
	return card, nil
}
