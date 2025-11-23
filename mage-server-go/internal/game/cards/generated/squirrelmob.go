package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Squirrel Mob", NewSquirrelMob)
}

// NewSquirrelMob creates a Squirrel Mob
// {1}{G}{G} - CREATURE
func NewSquirrelMob(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Squirrel Mob")
	card.ManaCost = "{1}{G}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"SQUIRREL"}
	card.Power = "2"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewBoostEffect(xValue, xValue)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
