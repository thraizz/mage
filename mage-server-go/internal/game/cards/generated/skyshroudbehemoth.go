package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Skyshroud Behemoth", NewSkyshroudBehemoth)
}

// NewSkyshroudBehemoth creates a Skyshroud Behemoth
// {5}{G}{G} - CREATURE
func NewSkyshroudBehemoth(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Skyshroud Behemoth")
	card.ManaCost = "{5}{G}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"BEAST"}
	card.Power = "10"
	card.Toughness = "10"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
