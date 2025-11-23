package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Looming Hoverguard", NewLoomingHoverguard)
}

// NewLoomingHoverguard creates a Looming Hoverguard
// {4}{U}{U} - CREATURE
// Flying
func NewLoomingHoverguard(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Looming Hoverguard")
	card.ManaCost = "{4}{U}{U}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"DRONE"}
	card.Power = "3"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	return card, nil
}
