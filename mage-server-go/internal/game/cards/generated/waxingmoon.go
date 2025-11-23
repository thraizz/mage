package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Waxing Moon", NewWaxingMoon)
}

// NewWaxingMoon creates a Waxing Moon
// {1}{G} - INSTANT
func NewWaxingMoon(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Waxing Moon")
	card.ManaCost = "{1}{G}"
	card.Types = []string{"INSTANT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - TransformTargetEffect()
	// card.AddAbility(ability0)
	return card, nil
}
