package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/effects"
)

func init() {
	cards.Register("Bronze Guardian", NewBronzeGuardian)
}

// NewBronzeGuardian creates a Bronze Guardian
// {4}{W} - ARTIFACT CREATURE
// DoubleStrike
func NewBronzeGuardian(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Bronze Guardian")
	card.ManaCost = "{4}{W}"
	card.Types = []string{"ARTIFACT", "CREATURE"}
	card.Subtypes = []string{"GOLEM"}
	card.Power = "0"
	card.Toughness = "5"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordDoubleStrike)
	card.AddAbility(ability0)
	ability1, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewGrantAbilityEffect(new WardAbility(), true)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability1)
	return card, nil
}