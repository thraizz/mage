package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Spellbreaker Behemoth", NewSpellbreakerBehemoth)
}

// NewSpellbreakerBehemoth creates a Spellbreaker Behemoth
// {1}{R}{G}{G} - CREATURE
func NewSpellbreakerBehemoth(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Spellbreaker Behemoth")
	card.ManaCost = "{1}{R}{G}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"BEAST"}
	card.Power = "5"
	card.Toughness = "5"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
