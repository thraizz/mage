package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Ghastlord Of Fugue", NewGhastlordOfFugue)
}

// NewGhastlordOfFugue creates a Ghastlord Of Fugue
// {U/B}{U/B}{U/B}{U/B}{U/B} - CREATURE
func NewGhastlordOfFugue(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Ghastlord Of Fugue")
	card.ManaCost = "{U/B}{U/B}{U/B}{U/B}{U/B}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"SPIRIT", "AVATAR"}
	card.Power = "4"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
