package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Containment Construct", NewContainmentConstruct)
}

// NewContainmentConstruct creates a Containment Construct
// {2} - ARTIFACT CREATURE
func NewContainmentConstruct(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Containment Construct")
	card.ManaCost = "{2}"
	card.Types = []string{"ARTIFACT", "CREATURE"}
	card.Subtypes = []string{"CONSTRUCT"}
	card.Power = "2"
	card.Toughness = "1"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
