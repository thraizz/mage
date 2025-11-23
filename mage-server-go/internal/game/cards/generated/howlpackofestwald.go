package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Howlpack Of Estwald", NewHowlpackOfEstwald)
}

// NewHowlpackOfEstwald creates a Howlpack Of Estwald
//   - CREATURE
func NewHowlpackOfEstwald(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Howlpack Of Estwald")
	card.ManaCost = ""
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"WEREWOLF"}
	card.Power = "4"
	card.Toughness = "6"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
