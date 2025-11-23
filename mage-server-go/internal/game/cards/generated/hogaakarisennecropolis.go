package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Hogaak Arisen Necropolis", NewHogaakArisenNecropolis)
}

// NewHogaakArisenNecropolis creates a Hogaak Arisen Necropolis
// {5}{B/G}{B/G} - CREATURE
// Trample
func NewHogaakArisenNecropolis(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Hogaak Arisen Necropolis")
	card.ManaCost = "{5}{B/G}{B/G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"AVATAR"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "8"
	card.Toughness = "8"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordTrample)
	card.AddAbility(ability0)
	return card, nil
}
