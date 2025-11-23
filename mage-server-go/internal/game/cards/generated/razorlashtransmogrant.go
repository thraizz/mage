package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Razorlash Transmogrant", NewRazorlashTransmogrant)
}

// NewRazorlashTransmogrant creates a Razorlash Transmogrant
// {2} - ARTIFACT CREATURE
func NewRazorlashTransmogrant(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Razorlash Transmogrant")
	card.ManaCost = "{2}"
	card.Types = []string{"ARTIFACT", "CREATURE"}
	card.Subtypes = []string{"ZOMBIE"}
	card.Power = "3"
	card.Toughness = "1"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
