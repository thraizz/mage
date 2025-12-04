package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Wake The Dead", NewWakeTheDead)
}

// NewWakeTheDead creates a Wake The Dead
// {X}{B}{B} - INSTANT
func NewWakeTheDead(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Wake The Dead")
	card.ManaCost = "{X}{B}{B}"
	card.Types = []string{"INSTANT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - SacrificeTargetEffect("Sacrifice those creatures at the beginning of the...)
	// card.AddAbility(ability0)
	return card, nil
}
