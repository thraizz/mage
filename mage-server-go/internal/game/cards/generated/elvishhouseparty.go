package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Elvish House Party", NewElvishHouseParty)
}

// NewElvishHouseParty creates a Elvish House Party
//
//	-
func NewElvishHouseParty(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Elvish House Party")
	card.ManaCost = ""
	card.Subtypes = []string{"ELF", "ROGUE"}
	card.Power = "0"
	card.Toughness = "0"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
