package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Cragsmasher Yeti", NewCragsmasherYeti)
}

// NewCragsmasherYeti creates a Cragsmasher Yeti
// {4}{R}{R} - CREATURE
// Trample
func NewCragsmasherYeti(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Cragsmasher Yeti")
	card.ManaCost = "{4}{R}{R}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"YETI"}
	card.Power = "4"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordTrample)
	card.AddAbility(ability0)
	return card, nil
}
