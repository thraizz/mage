package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Vadrok Apex Of Thunder", NewVadrokApexOfThunder)
}

// NewVadrokApexOfThunder creates a Vadrok Apex Of Thunder
// {U}{R}{W} - CREATURE
// Flying, FirstStrike
func NewVadrokApexOfThunder(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Vadrok Apex Of Thunder")
	card.ManaCost = "{U}{R}{W}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"ELEMENTAL", "DINOSAUR", "CAT"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "3"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	ability1 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFirstStrike)
	card.AddAbility(ability1)
	return card, nil
}
