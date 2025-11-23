package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Leyline Of Mutation", NewLeylineOfMutation)
}

// NewLeylineOfMutation creates a Leyline Of Mutation
// {2}{G}{G} - ENCHANTMENT
func NewLeylineOfMutation(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Leyline Of Mutation")
	card.ManaCost = "{2}{G}{G}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
