package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Malanthrope", NewMalanthrope)
}

// NewMalanthrope creates a Malanthrope
// {1}{G}{U} - CREATURE
// Flying
func NewMalanthrope(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Malanthrope")
	card.ManaCost = "{1}{G}{U}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"TYRANID"}
	card.Power = "2"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	return card, nil
}
