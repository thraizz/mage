package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Kylox Visionary Inventor", NewKyloxVisionaryInventor)
}

// NewKyloxVisionaryInventor creates a Kylox Visionary Inventor
// {5}{U}{R} - CREATURE
// Haste
func NewKyloxVisionaryInventor(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Kylox Visionary Inventor")
	card.ManaCost = "{5}{U}{R}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"LIZARD", "ARTIFICER"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "4"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordHaste)
	card.AddAbility(ability0)
	return card, nil
}
