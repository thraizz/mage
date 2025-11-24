package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Boggart Birth Rite", NewBoggartBirthRite)
}

// NewBoggartBirthRite creates a Boggart Birth Rite
// {B} - KINDRED SORCERY
func NewBoggartBirthRite(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Boggart Birth Rite")
	card.ManaCost = "{B}"
	card.Types = []string{"KINDRED", "SORCERY"}
	card.Subtypes = []string{"GOBLIN"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewReturnFromGraveyardToHandTargetEffect()).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}