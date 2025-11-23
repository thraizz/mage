package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Neyam Shai Murad", NewNeyamShaiMurad)
}

// NewNeyamShaiMurad creates a Neyam Shai Murad
// {2}{W}{B} - CREATURE
func NewNeyamShaiMurad(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Neyam Shai Murad")
	card.ManaCost = "{2}{W}{B}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"HUMAN", "ROGUE"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "3"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
