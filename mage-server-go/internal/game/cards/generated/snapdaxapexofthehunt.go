package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Snapdax Apex Of The Hunt", NewSnapdaxApexOfTheHunt)
}

// NewSnapdaxApexOfTheHunt creates a Snapdax Apex Of The Hunt
// {1}{R}{W}{B} - CREATURE
// DoubleStrike
func NewSnapdaxApexOfTheHunt(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Snapdax Apex Of The Hunt")
	card.ManaCost = "{1}{R}{W}{B}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"DINOSAUR", "CAT", "NIGHTMARE"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "3"
	card.Toughness = "5"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordDoubleStrike)
	card.AddAbility(ability0)
	return card, nil
}
