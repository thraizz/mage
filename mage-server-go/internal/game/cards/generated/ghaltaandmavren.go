package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Ghalta And Mavren", NewGhaltaAndMavren)
}

// NewGhaltaAndMavren creates a Ghalta And Mavren
// {3}{G}{G}{W}{W} - CREATURE
// Trample
func NewGhaltaAndMavren(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Ghalta And Mavren")
	card.ManaCost = "{3}{G}{G}{W}{W}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"DINOSAUR", "VAMPIRE"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "12"
	card.Toughness = "12"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordTrample)
	card.AddAbility(ability0)
	return card, nil
}
