package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Eiganjo Seat Of The Empire", NewEiganjoSeatOfTheEmpire)
}

// NewEiganjoSeatOfTheEmpire creates a Eiganjo Seat Of The Empire
//   - LAND
func NewEiganjoSeatOfTheEmpire(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Eiganjo Seat Of The Empire")
	card.ManaCost = ""
	card.Types = []string{"LAND"}
	card.Supertypes = []string{"LEGENDARY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.BuildSimpleManaAbility(card.ID, "W")
	card.AddAbility(ability0)
	return card, nil
}
