package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Incendiary Sabotage", NewIncendiarySabotage)
}

// NewIncendiarySabotage creates a Incendiary Sabotage
// {2}{R}{R} - INSTANT
func NewIncendiarySabotage(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Incendiary Sabotage")
	card.ManaCost = "{2}{R}{R}"
	card.Types = []string{"INSTANT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - DamageAllEffect(3, new FilterCreaturePermanent())
	// card.AddAbility(ability0)
	return card, nil
}
