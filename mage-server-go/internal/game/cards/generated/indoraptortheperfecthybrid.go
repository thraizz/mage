package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Indoraptor The Perfect Hybrid", NewIndoraptorThePerfectHybrid)
}

// NewIndoraptorThePerfectHybrid creates a Indoraptor The Perfect Hybrid
// {1}{B/G}{R} - CREATURE
func NewIndoraptorThePerfectHybrid(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Indoraptor The Perfect Hybrid")
	card.ManaCost = "{1}{B/G}{R}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"DINOSAUR", "MUTANT"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "3"
	card.Toughness = "1"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}