package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Spatial Binding", NewSpatialBinding)
}

// NewSpatialBinding creates a Spatial Binding
// {U}{B} - ENCHANTMENT
func NewSpatialBinding(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Spatial Binding")
	card.ManaCost = "{U}{B}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement activated ability with unmapped effects
	//   - SpatialBindingReplacementEffect()
	// card.AddAbility(ability0)
	return card, nil
}
