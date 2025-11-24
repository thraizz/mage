package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Beluna Grandsquall", NewBelunaGrandsquall)
}

// NewBelunaGrandsquall creates a Beluna Grandsquall
// {G}{U}{R} - CREATURE
// Trample
func NewBelunaGrandsquall(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Beluna Grandsquall")
	card.ManaCost = "{G}{U}{R}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"GIANT", "NOBLE"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "4"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordTrample)
	card.AddAbility(ability0)
	return card, nil
}