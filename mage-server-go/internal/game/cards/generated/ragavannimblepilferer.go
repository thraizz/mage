package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Ragavan Nimble Pilferer", NewRagavanNimblePilferer)
}

// NewRagavanNimblePilferer creates a Ragavan Nimble Pilferer
// {R} - CREATURE
func NewRagavanNimblePilferer(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Ragavan Nimble Pilferer")
	card.ManaCost = "{R}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"MONKEY", "PIRATE"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "2"
	card.Toughness = "1"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}