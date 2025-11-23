package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Yidaro Wandering Monster", NewYidaroWanderingMonster)
}

// NewYidaroWanderingMonster creates a Yidaro Wandering Monster
// {5}{R}{R} - CREATURE
// Trample, Haste
func NewYidaroWanderingMonster(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Yidaro Wandering Monster")
	card.ManaCost = "{5}{R}{R}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"DINOSAUR", "TURTLE"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "8"
	card.Toughness = "8"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordTrample)
	card.AddAbility(ability0)
	ability1 := abilities.NewKeywordAbility(card.ID, abilities.KeywordHaste)
	card.AddAbility(ability1)
	return card, nil
}
