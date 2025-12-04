package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Lorehold Command", NewLoreholdCommand)
}

// NewLoreholdCommand creates a Lorehold Command
// {3}{R}{W} - INSTANT
func NewLoreholdCommand(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Lorehold Command")
	card.ManaCost = "{3}{R}{W}"
	card.Types = []string{"INSTANT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - SacrificeControllerEffect(StaticFilters.FILTER_PERMANENT, 1, "")
	// card.AddAbility(ability0)
	return card, nil
}
