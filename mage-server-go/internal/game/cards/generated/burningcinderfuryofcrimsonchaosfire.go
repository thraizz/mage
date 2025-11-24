package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Burning Cinder Fury Of Crimson Chaos Fire", NewBurningCinderFuryOfCrimsonChaosFire)
}

// NewBurningCinderFuryOfCrimsonChaosFire creates a Burning Cinder Fury Of Crimson Chaos Fire
// {3}{R} - ENCHANTMENT
func NewBurningCinderFuryOfCrimsonChaosFire(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Burning Cinder Fury Of Crimson Chaos Fire")
	card.ManaCost = "{3}{R}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewDamageEffect(3)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}