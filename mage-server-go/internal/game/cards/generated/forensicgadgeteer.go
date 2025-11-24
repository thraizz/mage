package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Forensic Gadgeteer", NewForensicGadgeteer)
}

// NewForensicGadgeteer creates a Forensic Gadgeteer
// {2}{U} - CREATURE
func NewForensicGadgeteer(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Forensic Gadgeteer")
	card.ManaCost = "{2}{U}"
	card.Types = []string{"CREATURE"}
	card.Power = "2"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
