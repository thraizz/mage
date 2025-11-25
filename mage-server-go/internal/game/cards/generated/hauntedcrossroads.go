package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Haunted Crossroads", NewHauntedCrossroads)
}

// NewHauntedCrossroads creates a Haunted Crossroads
// {2}{B} - ENCHANTMENT
func NewHauntedCrossroads(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Haunted Crossroads")
	card.ManaCost = "{2}{B}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement activated ability with unmapped effects
	//   - PutOnLibraryTargetEffect()
	// card.AddAbility(ability0)
	return card, nil
}
