package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Bilbo Birthday Celebrant", NewBilboBirthdayCelebrant)
}

// NewBilboBirthdayCelebrant creates a Bilbo Birthday Celebrant
// {W}{B}{G} - CREATURE
func NewBilboBirthdayCelebrant(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Bilbo Birthday Celebrant")
	card.ManaCost = "{W}{B}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"HALFLING", "ROGUE"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "2"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
