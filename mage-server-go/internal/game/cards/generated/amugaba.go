package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Amugaba", NewAmugaba)
}

// NewAmugaba creates a Amugaba
// {5}{U}{U} - CREATURE
// Flying
func NewAmugaba(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Amugaba")
	card.ManaCost = "{5}{U}{U}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"ILLUSION"}
	card.Power = "6"
	card.Toughness = "6"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	ability1 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddEffect(abilities.NewReturnToHandSourceEffect()).
		Build()
	card.AddAbility(ability1)
	return card, nil
}
