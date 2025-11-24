package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Bogbrew Witch", NewBogbrewWitch)
}

// NewBogbrewWitch creates a Bogbrew Witch
// {3}{B} - CREATURE
func NewBogbrewWitch(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Bogbrew Witch")
	card.ManaCost = "{3}{B}"
	card.Types = []string{"CREATURE"}
	card.Power = "1"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddTapCost().
		// TODO: SearchLibraryPutInPlayEffect with complex parameters
		Build()
	card.AddAbility(ability0)
	return card, nil
}
