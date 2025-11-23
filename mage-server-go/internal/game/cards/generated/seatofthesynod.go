package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Seat Of The Synod", NewSeatOfTheSynod)
}

// NewSeatOfTheSynod creates a Seat Of The Synod
//   - ARTIFACT LAND
func NewSeatOfTheSynod(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Seat Of The Synod")
	card.ManaCost = ""
	card.Types = []string{"ARTIFACT", "LAND"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.BuildSimpleManaAbility(card.ID, "U")
	card.AddAbility(ability0)
	return card, nil
}
