package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("The Gitrog Monster", NewTheGitrogMonster)
}

// NewTheGitrogMonster creates a The Gitrog Monster
// {3}{B}{G} - CREATURE
// Deathtouch
func NewTheGitrogMonster(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "The Gitrog Monster")
	card.ManaCost = "{3}{B}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"FROG", "HORROR"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "6"
	card.Toughness = "6"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordDeathtouch)
	card.AddAbility(ability0)
	return card, nil
}
