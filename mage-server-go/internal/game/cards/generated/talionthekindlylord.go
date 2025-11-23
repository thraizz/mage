package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Talion The Kindly Lord", NewTalionTheKindlyLord)
}

// NewTalionTheKindlyLord creates a Talion The Kindly Lord
// {2}{U}{B} - CREATURE
// Flying
func NewTalionTheKindlyLord(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Talion The Kindly Lord")
	card.ManaCost = "{2}{U}{B}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"FAERIE", "NOBLE"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "3"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	return card, nil
}
