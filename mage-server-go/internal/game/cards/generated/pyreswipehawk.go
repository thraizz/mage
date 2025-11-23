package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Pyreswipe Hawk", NewPyreswipeHawk)
}

// NewPyreswipeHawk creates a Pyreswipe Hawk
// {3}{R}{R} - CREATURE
// Flying, Haste
func NewPyreswipeHawk(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Pyreswipe Hawk")
	card.ManaCost = "{3}{R}{R}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"ELEMENTAL", "BIRD"}
	card.Power = "4"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	ability1 := abilities.NewKeywordAbility(card.ID, abilities.KeywordHaste)
	card.AddAbility(ability1)
	ability2, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewBoostEffect(GreatestAmongPermanentsValue.MANAVALUE_CONTROLLED_ARTIFACTS, StaticValue.get(0))).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability2)
	return card, nil
}
