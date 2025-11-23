package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Ayumi The Last Visitor", NewAyumiTheLastVisitor)
}

// NewAyumiTheLastVisitor creates a Ayumi The Last Visitor
// {3}{G}{G} - CREATURE
func NewAyumiTheLastVisitor(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Ayumi The Last Visitor")
	card.ManaCost = "{3}{G}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"SPIRIT"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "7"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
