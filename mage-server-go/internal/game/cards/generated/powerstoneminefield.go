package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Powerstone Minefield", NewPowerstoneMinefield)
}

// NewPowerstoneMinefield creates a Powerstone Minefield
// {2}{R}{W} - ENCHANTMENT
func NewPowerstoneMinefield(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Powerstone Minefield")
	card.ManaCost = "{2}{R}{W}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
