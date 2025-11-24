package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Symbiote Spider Man", NewSymbioteSpiderMan)
}

// NewSymbioteSpiderMan creates a Symbiote Spider Man
// {2}{U/B} - CREATURE
func NewSymbioteSpiderMan(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Symbiote Spider Man")
	card.ManaCost = "{2}{U/B}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"SYMBIOTE", "SPIDER", "HERO"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "2"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}