package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Breaching Leviathan", NewBreachingLeviathan)
}

// NewBreachingLeviathan creates a Breaching Leviathan
// {7}{U}{U} - CREATURE
func NewBreachingLeviathan(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Breaching Leviathan")
	card.ManaCost = "{7}{U}{U}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"LEVIATHAN"}
	card.Power = "9"
	card.Toughness = "9"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
