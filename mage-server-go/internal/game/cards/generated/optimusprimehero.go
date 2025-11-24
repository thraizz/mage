package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Optimus Prime Hero", NewOptimusPrimeHero)
}

// NewOptimusPrimeHero creates a Optimus Prime Hero
// {3}{U}{R}{W} - ARTIFACT CREATURE
func NewOptimusPrimeHero(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Optimus Prime Hero")
	card.ManaCost = "{3}{U}{R}{W}"
	card.Types = []string{"ARTIFACT", "CREATURE"}
	card.Subtypes = []string{"ROBOT"}
	card.Power = "4"
	card.Toughness = "8"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
