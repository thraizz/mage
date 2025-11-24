package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Go Shintai Of Boundless Vigor", NewGoShintaiOfBoundlessVigor)
}

// NewGoShintaiOfBoundlessVigor creates a Go Shintai Of Boundless Vigor
// {1}{G} - ENCHANTMENT CREATURE
// Trample
func NewGoShintaiOfBoundlessVigor(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Go Shintai Of Boundless Vigor")
	card.ManaCost = "{1}{G}"
	card.Types = []string{"ENCHANTMENT", "CREATURE"}
	card.Subtypes = []string{"SHRINE"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "1"
	card.Toughness = "1"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordTrample)
	card.AddAbility(ability0)
	return card, nil
}
