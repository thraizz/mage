package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Thassa Deep Dwelling", NewThassaDeepDwelling)
}

// NewThassaDeepDwelling creates a Thassa Deep Dwelling
// {3}{U} - ENCHANTMENT CREATURE
// Indestructible
func NewThassaDeepDwelling(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Thassa Deep Dwelling")
	card.ManaCost = "{3}{U}"
	card.Types = []string{"ENCHANTMENT", "CREATURE"}
	card.Subtypes = []string{"GOD"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "6"
	card.Toughness = "5"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordIndestructible)
	card.AddAbility(ability0)
	ability1 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddEffect(abilities.NewTapEffect("tap another target creature")).
		Build()
	card.AddAbility(ability1)
	return card, nil
}
