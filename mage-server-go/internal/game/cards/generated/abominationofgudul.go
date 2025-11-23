package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Abomination Of Gudul", NewAbominationOfGudul)
}

// NewAbominationOfGudul creates a Abomination Of Gudul
// {3}{B}{G}{U} - CREATURE
// Flying
func NewAbominationOfGudul(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Abomination Of Gudul")
	card.ManaCost = "{3}{B}{G}{U}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"HORROR"}
	card.Power = "3"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	return card, nil
}
