package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Kaya Geist Hunter", NewKayaGeistHunter)
}

// NewKayaGeistHunter creates a Kaya Geist Hunter
// {1}{W}{B} - PLANESWALKER
func NewKayaGeistHunter(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Kaya Geist Hunter")
	card.ManaCost = "{1}{W}{B}"
	card.Types = []string{"PLANESWALKER"}
	card.Subtypes = []string{"KAYA"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Loyalty = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
