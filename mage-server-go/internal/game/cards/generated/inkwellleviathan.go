package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Inkwell Leviathan", NewInkwellLeviathan)
}

// NewInkwellLeviathan creates a Inkwell Leviathan
// {7}{U}{U} - ARTIFACT CREATURE
// Trample
func NewInkwellLeviathan(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Inkwell Leviathan")
	card.ManaCost = "{7}{U}{U}"
	card.Types = []string{"ARTIFACT", "CREATURE"}
	card.Subtypes = []string{"LEVIATHAN"}
	card.Power = "7"
	card.Toughness = "11"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordTrample)
	card.AddAbility(ability0)
	return card, nil
}
