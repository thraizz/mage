package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Frenetic Efreet", NewFreneticEfreet)
}

// NewFreneticEfreet creates a Frenetic Efreet
// {1}{U}{R} - CREATURE
// Flying
func NewFreneticEfreet(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Frenetic Efreet")
	card.ManaCost = "{1}{U}{R}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"EFREET"}
	card.Power = "2"
	card.Toughness = "1"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	return card, nil
}
