package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Seer Of Stolen Sight", NewSeerOfStolenSight)
}

// NewSeerOfStolenSight creates a Seer Of Stolen Sight
// {2}{B} - CREATURE
// Menace
func NewSeerOfStolenSight(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Seer Of Stolen Sight")
	card.ManaCost = "{2}{B}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"PHYREXIAN", "WARLOCK"}
	card.Power = "2"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordMenace)
	card.AddAbility(ability0)
	return card, nil
}
