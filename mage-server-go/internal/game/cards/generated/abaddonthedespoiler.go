package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Abaddon The Despoiler", NewAbaddonTheDespoiler)
}

// NewAbaddonTheDespoiler creates a Abaddon The Despoiler
// {2}{U}{B}{R} - CREATURE
// Trample
func NewAbaddonTheDespoiler(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Abaddon The Despoiler")
	card.ManaCost = "{2}{U}{B}{R}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"ASTARTES", "WARRIOR"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "5"
	card.Toughness = "5"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordTrample)
	card.AddAbility(ability0)
	return card, nil
}
