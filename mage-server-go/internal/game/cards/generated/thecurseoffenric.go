package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("The Curse Of Fenric", NewTheCurseOfFenric)
}

// NewTheCurseOfFenric creates a The Curse Of Fenric
// {2}{G}{W} - ENCHANTMENT
func NewTheCurseOfFenric(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "The Curse Of Fenric")
	card.ManaCost = "{2}{G}{W}"
	card.Types = []string{"ENCHANTMENT"}
	card.Subtypes = []string{"SAGA", "HORROR"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "6"
	card.Toughness = "6"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
