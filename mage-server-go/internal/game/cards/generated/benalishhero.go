package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Benalish Hero", NewBenalishHero)
}

// NewBenalishHero creates a Benalish Hero
// {W} - CREATURE
func NewBenalishHero(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Benalish Hero")
	card.ManaCost = "{W}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"HUMAN", "SOLDIER"}
	card.Power = "1"
	card.Toughness = "1"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
