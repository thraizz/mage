package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Birgi God Of Storytelling", NewBirgiGodOfStorytelling)
}

// NewBirgiGodOfStorytelling creates a Birgi God Of Storytelling
//   - CREATURE
func NewBirgiGodOfStorytelling(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Birgi God Of Storytelling")
	card.ManaCost = ""
	card.Types = []string{"CREATURE"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
