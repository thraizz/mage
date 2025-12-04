package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Ravens Run Dragoon", NewRavensRunDragoon)
}

// NewRavensRunDragoon creates a Ravens Run Dragoon
// {2}{G/W}{G/W} - CREATURE
func NewRavensRunDragoon(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Ravens Run Dragoon")
	card.ManaCost = "{2}{G/W}{G/W}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"ELF", "KNIGHT"}
	card.Power = "3"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
