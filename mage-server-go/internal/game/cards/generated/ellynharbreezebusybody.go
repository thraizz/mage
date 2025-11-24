package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Ellyn Harbreeze Busybody", NewEllynHarbreezeBusybody)
}

// NewEllynHarbreezeBusybody creates a Ellyn Harbreeze Busybody
// {3}{W} - CREATURE
func NewEllynHarbreezeBusybody(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Ellyn Harbreeze Busybody")
	card.ManaCost = "{3}{W}"
	card.Types = []string{"CREATURE"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "2"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}