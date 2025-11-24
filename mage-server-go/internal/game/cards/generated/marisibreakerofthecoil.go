package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Marisi Breaker Of The Coil", NewMarisiBreakerOfTheCoil)
}

// NewMarisiBreakerOfTheCoil creates a Marisi Breaker Of The Coil
// {1}{R}{G}{W} - CREATURE
func NewMarisiBreakerOfTheCoil(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Marisi Breaker Of The Coil")
	card.ManaCost = "{1}{R}{G}{W}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"CAT", "WARRIOR"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "5"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
