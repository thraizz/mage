package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Beast Kin Ranger", NewBeastKinRanger)
}

// NewBeastKinRanger creates a Beast Kin Ranger
// {2}{G} - CREATURE
// Trample
func NewBeastKinRanger(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Beast Kin Ranger")
	card.ManaCost = "{2}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"ELF", "RANGER"}
	card.Power = "3"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordTrample)
	card.AddAbility(ability0)
	ability1, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewBoostEffect(1, 0)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability1)
	return card, nil
}
