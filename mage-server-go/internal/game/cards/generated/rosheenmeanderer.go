package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Rosheen Meanderer", NewRosheenMeanderer)
}

// NewRosheenMeanderer creates a Rosheen Meanderer
// {3}{R/G} - CREATURE
func NewRosheenMeanderer(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Rosheen Meanderer")
	card.ManaCost = "{3}{R/G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"GIANT", "SHAMAN"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "4"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
