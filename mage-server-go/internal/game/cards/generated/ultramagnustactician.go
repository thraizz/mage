package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Ultra Magnus Tactician", NewUltraMagnusTactician)
}

// NewUltraMagnusTactician creates a Ultra Magnus Tactician
// {4}{R}{G}{W} - ARTIFACT CREATURE
func NewUltraMagnusTactician(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Ultra Magnus Tactician")
	card.ManaCost = "{4}{R}{G}{W}"
	card.Types = []string{"ARTIFACT", "CREATURE"}
	card.Subtypes = []string{"ROBOT"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "7"
	card.Toughness = "7"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
