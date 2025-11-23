package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Conquerors Galleon", NewConquerorsGalleon)
}

// NewConquerorsGalleon creates a Conquerors Galleon
// {4} - ARTIFACT
func NewConquerorsGalleon(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Conquerors Galleon")
	card.ManaCost = "{4}"
	card.Types = []string{"ARTIFACT"}
	card.Subtypes = []string{"VEHICLE"}
	card.Power = "2"
	card.Toughness = "10"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
