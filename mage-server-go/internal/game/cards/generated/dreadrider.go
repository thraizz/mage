package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Dread Rider", NewDreadRider)
}

// NewDreadRider creates a Dread Rider
// {5}{B} - CREATURE
func NewDreadRider(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Dread Rider")
	card.ManaCost = "{5}{B}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"SPIRIT", "KNIGHT"}
	card.Power = "3"
	card.Toughness = "7"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddTapCost().
		AddEffect(abilities.NewLoseLifeEffect(3)).
		Build()
	card.AddAbility(ability0)
	return card, nil
}
