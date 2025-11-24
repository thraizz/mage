package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Stonewing Antagonizer", NewStonewingAntagonizer)
}

// NewStonewingAntagonizer creates a Stonewing Antagonizer
//   - ARTIFACT CREATURE
//
// Flying
func NewStonewingAntagonizer(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Stonewing Antagonizer")
	card.ManaCost = ""
	card.Types = []string{"ARTIFACT", "CREATURE"}
	card.Subtypes = []string{"GARGOYLE", "HORROR"}
	card.Power = "4"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	return card, nil
}
