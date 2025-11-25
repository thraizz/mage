package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Walk In Closet Forgotten Cellar", NewWalkInClosetForgottenCellar)
}

// NewWalkInClosetForgottenCellar creates a Walk In Closet Forgotten Cellar
// {2}{G} - ENCHANTMENT
func NewWalkInClosetForgottenCellar(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Walk In Closet Forgotten Cellar")
	card.ManaCost = "{2}{G}"
	card.Types = []string{"ENCHANTMENT"}
	card.Subtypes = []string{"ROOM"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
