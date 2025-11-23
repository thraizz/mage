package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("The Sixth Doctor", NewTheSixthDoctor)
}

// NewTheSixthDoctor creates a The Sixth Doctor
// {4}{G}{U} - CREATURE
func NewTheSixthDoctor(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "The Sixth Doctor")
	card.ManaCost = "{4}{G}{U}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"TIME_LORD", "DOCTOR"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "3"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
