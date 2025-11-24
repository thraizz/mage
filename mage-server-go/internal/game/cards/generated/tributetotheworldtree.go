package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Tribute To The World Tree", NewTributeToTheWorldTree)
}

// NewTributeToTheWorldTree creates a Tribute To The World Tree
// {G}{G}{G} - ENCHANTMENT
func NewTributeToTheWorldTree(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Tribute To The World Tree")
	card.ManaCost = "{G}{G}{G}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}