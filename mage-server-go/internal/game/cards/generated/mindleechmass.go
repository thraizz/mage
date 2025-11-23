package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Mindleech Mass", NewMindleechMass)
}

// NewMindleechMass creates a Mindleech Mass
// {5}{U}{B}{B} - CREATURE
// Trample
func NewMindleechMass(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Mindleech Mass")
	card.ManaCost = "{5}{U}{B}{B}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"HORROR"}
	card.Power = "6"
	card.Toughness = "6"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordTrample)
	card.AddAbility(ability0)
	return card, nil
}
