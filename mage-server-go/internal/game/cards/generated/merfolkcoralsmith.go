package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Merfolk Coralsmith", NewMerfolkCoralsmith)
}

// NewMerfolkCoralsmith creates a Merfolk Coralsmith
// {2}{U} - CREATURE
func NewMerfolkCoralsmith(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Merfolk Coralsmith")
	card.ManaCost = "{2}{U}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"MERFOLK"}
	card.Power = "2"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewScryEffect(1)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
