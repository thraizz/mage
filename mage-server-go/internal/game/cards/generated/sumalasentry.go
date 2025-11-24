package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Sumala Sentry", NewSumalaSentry)
}

// NewSumalaSentry creates a Sumala Sentry
// {G}{W} - CREATURE
// Reach
func NewSumalaSentry(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Sumala Sentry")
	card.ManaCost = "{G}{W}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"ELF", "ARCHER"}
	card.Power = "1"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordReach)
	card.AddAbility(ability0)
	return card, nil
}
