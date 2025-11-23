package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Dangerous Wager", NewDangerousWager)
}

// NewDangerousWager creates a Dangerous Wager
// {1}{R} - INSTANT
func NewDangerousWager(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Dangerous Wager")
	card.ManaCost = "{1}{R}"
	card.Types = []string{"INSTANT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - DiscardHandControllerEffect()
	// card.AddAbility(ability0)
	return card, nil
}
