package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Icebreaker Kraken", NewIcebreakerKraken)
}

// NewIcebreakerKraken creates a Icebreaker Kraken
// {10}{U}{U} - CREATURE
func NewIcebreakerKraken(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Icebreaker Kraken")
	card.ManaCost = "{10}{U}{U}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"KRAKEN"}
	card.Supertypes = []string{"SNOW"}
	card.Power = "8"
	card.Toughness = "8"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
