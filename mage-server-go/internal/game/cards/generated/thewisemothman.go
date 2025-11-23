package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("The Wise Mothman", NewTheWiseMothman)
}

// NewTheWiseMothman creates a The Wise Mothman
// {1}{B}{G}{U} - CREATURE
// Flying
func NewTheWiseMothman(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "The Wise Mothman")
	card.ManaCost = "{1}{B}{G}{U}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"INSECT", "MUTANT"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "3"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	return card, nil
}
