package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Fanatical Devotion", NewFanaticalDevotion)
}

// NewFanaticalDevotion creates a Fanatical Devotion
// {2}{W} - ENCHANTMENT
func NewFanaticalDevotion(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Fanatical Devotion")
	card.ManaCost = "{2}{W}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement activated ability with unmapped effects
	//   - RegenerateTargetEffect()
	// card.AddAbility(ability0)
	return card, nil
}
