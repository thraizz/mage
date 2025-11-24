package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Summoning Trap", NewSummoningTrap)
}

// NewSummoningTrap creates a Summoning Trap
// {4}{G}{G} - INSTANT
func NewSummoningTrap(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Summoning Trap")
	card.ManaCost = "{4}{G}{G}"
	card.Types = []string{"INSTANT"}
	card.Subtypes = []string{"TRAP"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
