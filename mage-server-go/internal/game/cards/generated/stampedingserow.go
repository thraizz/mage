package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Stampeding Serow", NewStampedingSerow)
}

// NewStampedingSerow creates a Stampeding Serow
// {2}{G}{G} - CREATURE
// Trample
func NewStampedingSerow(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Stampeding Serow")
	card.ManaCost = "{2}{G}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"ANTELOPE", "BEAST"}
	card.Power = "5"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordTrample)
	card.AddAbility(ability0)
	return card, nil
}
