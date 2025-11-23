package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/effects"
)

func init() {
	cards.Register("Soaring Lightbringer", NewSoaringLightbringer)
}

// NewSoaringLightbringer creates a Soaring Lightbringer
// {4}{W} - ENCHANTMENT CREATURE
// Flying
func NewSoaringLightbringer(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Soaring Lightbringer")
	card.ManaCost = "{4}{W}"
	card.Types = []string{"ENCHANTMENT", "CREATURE"}
	card.Subtypes = []string{"BIRD", "GLIMMER"}
	card.Power = "4"
	card.Toughness = "5"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	ability1, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewGrantAbilityEffect("FlyingAbility", effects.DurationWhileOnBattlefield)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability1)
	return card, nil
}
