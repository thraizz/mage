package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Jasmine Boreal Of The Seven", NewJasmineBorealOfTheSeven)
}

// NewJasmineBorealOfTheSeven creates a Jasmine Boreal Of The Seven
// {1}{G}{W} - CREATURE
func NewJasmineBorealOfTheSeven(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Jasmine Boreal Of The Seven")
	card.ManaCost = "{1}{G}{W}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"HUMAN", "DRUID"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "2"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
