package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Chromium The Mutable", NewChromiumTheMutable)
}

// NewChromiumTheMutable creates a Chromium The Mutable
// {4}{W}{U}{B} - CREATURE
// Flash, Flying, Hexproof
func NewChromiumTheMutable(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Chromium The Mutable")
	card.ManaCost = "{4}{W}{U}{B}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"ELDER", "DRAGON", "HUMAN"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "7"
	card.Toughness = "7"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlash)
	card.AddAbility(ability0)
	ability1 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability1)
	ability2 := abilities.NewKeywordAbility(card.ID, abilities.KeywordHexproof)
	card.AddAbility(ability2)
	return card, nil
}
