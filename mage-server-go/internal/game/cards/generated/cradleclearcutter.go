package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Cradle Clearcutter", NewCradleClearcutter)
}

// NewCradleClearcutter creates a Cradle Clearcutter
// {6} - ARTIFACT CREATURE
func NewCradleClearcutter(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Cradle Clearcutter")
	card.ManaCost = "{6}"
	card.Types = []string{"ARTIFACT", "CREATURE"}
	card.Subtypes = []string{"GOLEM"}
	card.Power = "3"
	card.Toughness = "6"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
