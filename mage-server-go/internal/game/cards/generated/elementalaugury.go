package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Elemental Augury", NewElementalAugury)
}

// NewElementalAugury creates a Elemental Augury
// {U}{B}{R} - ENCHANTMENT
func NewElementalAugury(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Elemental Augury")
	card.ManaCost = "{U}{B}{R}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement activated ability with unmapped effects
	//   - ElementalAuguryEffect()
	// card.AddAbility(ability0)
	return card, nil
}
