package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Xira Arien", NewXiraArien)
}

// NewXiraArien creates a Xira Arien
// {B}{R}{G} - CREATURE
// Flying
func NewXiraArien(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Xira Arien")
	card.ManaCost = "{B}{R}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"INSECT", "WIZARD"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "1"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	ability1 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddTapCost().
		AddEffect(abilities.NewDrawCardsEffect(1)).
		Build()
	card.AddAbility(ability1)
	return card, nil
}
