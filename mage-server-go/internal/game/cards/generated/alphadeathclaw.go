package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Alpha Deathclaw", NewAlphaDeathclaw)
}

// NewAlphaDeathclaw creates a Alpha Deathclaw
// {4}{B}{G} - CREATURE
// Trample
func NewAlphaDeathclaw(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Alpha Deathclaw")
	card.ManaCost = "{4}{B}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"LIZARD", "MUTANT"}
	card.Power = "6"
	card.Toughness = "6"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordTrample)
	card.AddAbility(ability0)
	return card, nil
}