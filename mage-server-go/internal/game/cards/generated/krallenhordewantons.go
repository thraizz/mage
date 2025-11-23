package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Krallenhorde Wantons", NewKrallenhordeWantons)
}

// NewKrallenhordeWantons creates a Krallenhorde Wantons
//   - CREATURE
func NewKrallenhordeWantons(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Krallenhorde Wantons")
	card.ManaCost = ""
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"WEREWOLF"}
	card.Power = "7"
	card.Toughness = "7"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
