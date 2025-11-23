package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Sanguine Syphoner", NewSanguineSyphoner)
}

// NewSanguineSyphoner creates a Sanguine Syphoner
// {1}{B} - CREATURE
func NewSanguineSyphoner(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Sanguine Syphoner")
	card.ManaCost = "{1}{B}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"VAMPIRE", "WARLOCK"}
	card.Power = "1"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
