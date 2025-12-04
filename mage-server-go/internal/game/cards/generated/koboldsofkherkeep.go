package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Kobolds Of Kher Keep", NewKoboldsOfKherKeep)
}

// NewKoboldsOfKherKeep creates a Kobolds Of Kher Keep
// {0} - CREATURE
func NewKoboldsOfKherKeep(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Kobolds Of Kher Keep")
	card.ManaCost = "{0}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"KOBOLD"}
	card.Power = "0"
	card.Toughness = "1"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
