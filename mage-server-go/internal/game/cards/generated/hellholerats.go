package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Hellhole Rats", NewHellholeRats)
}

// NewHellholeRats creates a Hellhole Rats
// {2}{B}{R} - CREATURE
// Haste
func NewHellholeRats(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Hellhole Rats")
	card.ManaCost = "{2}{B}{R}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"RAT"}
	card.Power = "2"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordHaste)
	card.AddAbility(ability0)
	return card, nil
}
