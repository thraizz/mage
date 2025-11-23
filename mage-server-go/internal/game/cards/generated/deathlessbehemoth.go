package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Deathless Behemoth", NewDeathlessBehemoth)
}

// NewDeathlessBehemoth creates a Deathless Behemoth
// {6} - CREATURE
// Vigilance
func NewDeathlessBehemoth(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Deathless Behemoth")
	card.ManaCost = "{6}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"ELDRAZI"}
	card.Power = "6"
	card.Toughness = "6"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordVigilance)
	card.AddAbility(ability0)
	return card, nil
}
