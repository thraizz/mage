package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Yamabushis Storm", NewYamabushisStorm)
}

// NewYamabushisStorm creates a Yamabushis Storm
// {1}{R} - SORCERY
func NewYamabushisStorm(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Yamabushis Storm")
	card.ManaCost = "{1}{R}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - DamageAllEffect(1, new FilterCreaturePermanent())
	// card.AddAbility(ability0)
	return card, nil
}
