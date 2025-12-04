package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Gallifrey Stands", NewGallifreyStands)
}

// NewGallifreyStands creates a Gallifrey Stands
// {4}{W}{U} - ENCHANTMENT
func NewGallifreyStands(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Gallifrey Stands")
	card.ManaCost = "{4}{W}{U}"
	card.Types = []string{"ENCHANTMENT"}
	card.Supertypes = []string{"LEGENDARY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
