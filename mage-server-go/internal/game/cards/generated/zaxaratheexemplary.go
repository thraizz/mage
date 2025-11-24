package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Zaxara The Exemplary", NewZaxaraTheExemplary)
}

// NewZaxaraTheExemplary creates a Zaxara The Exemplary
// {1}{B}{G}{U} - CREATURE
// Deathtouch
func NewZaxaraTheExemplary(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Zaxara The Exemplary")
	card.ManaCost = "{1}{B}{G}{U}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"NIGHTMARE", "HYDRA"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "2"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordDeathtouch)
	card.AddAbility(ability0)
	return card, nil
}