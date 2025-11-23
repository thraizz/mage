package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Erg Raiders", NewErgRaiders)
}

// NewErgRaiders creates a Erg Raiders
// {1}{B} - CREATURE
func NewErgRaiders(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Erg Raiders")
	card.ManaCost = "{1}{B}"
	card.Types = []string{"CREATURE"}
	card.Power = "2"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
