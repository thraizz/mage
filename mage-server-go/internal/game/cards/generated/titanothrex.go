package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Titanoth Rex", NewTitanothRex)
}

// NewTitanothRex creates a Titanoth Rex
// {7}{G}{G} - CREATURE
// Trample
func NewTitanothRex(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Titanoth Rex")
	card.ManaCost = "{7}{G}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"DINOSAUR", "BEAST"}
	card.Power = "11"
	card.Toughness = "11"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordTrample)
	card.AddAbility(ability0)
	return card, nil
}
