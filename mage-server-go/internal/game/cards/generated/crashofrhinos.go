package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Crash Of Rhinos", NewCrashOfRhinos)
}

// NewCrashOfRhinos creates a Crash Of Rhinos
// {6}{G}{G} - CREATURE
// Trample
func NewCrashOfRhinos(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Crash Of Rhinos")
	card.ManaCost = "{6}{G}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"RHINO"}
	card.Power = "8"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordTrample)
	card.AddAbility(ability0)
	return card, nil
}
