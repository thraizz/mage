package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Kynaios And Tiro Of Meletis", NewKynaiosAndTiroOfMeletis)
}

// NewKynaiosAndTiroOfMeletis creates a Kynaios And Tiro Of Meletis
// {R}{G}{W}{U} - CREATURE
func NewKynaiosAndTiroOfMeletis(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Kynaios And Tiro Of Meletis")
	card.ManaCost = "{R}{G}{W}{U}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"HUMAN", "SOLDIER"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "2"
	card.Toughness = "8"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
