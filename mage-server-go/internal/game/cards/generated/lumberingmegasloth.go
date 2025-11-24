package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Lumbering Megasloth", NewLumberingMegasloth)
}

// NewLumberingMegasloth creates a Lumbering Megasloth
// {10}{G}{G} - CREATURE
// Trample
func NewLumberingMegasloth(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Lumbering Megasloth")
	card.ManaCost = "{10}{G}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"SLOTH", "MUTANT"}
	card.Power = "8"
	card.Toughness = "8"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordTrample)
	card.AddAbility(ability0)
	return card, nil
}