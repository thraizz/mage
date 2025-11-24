package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("The Ur Dragon", NewTheUrDragon)
}

// NewTheUrDragon creates a The Ur Dragon
// {4}{W}{U}{B}{R}{G} - CREATURE
// Flying
func NewTheUrDragon(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "The Ur Dragon")
	card.ManaCost = "{4}{W}{U}{B}{R}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"DRAGON", "AVATAR"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "10"
	card.Toughness = "10"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	return card, nil
}
