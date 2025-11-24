package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Apostle Of Purifying Light", NewApostleOfPurifyingLight)
}

// NewApostleOfPurifyingLight creates a Apostle Of Purifying Light
// {1}{W} - CREATURE
func NewApostleOfPurifyingLight(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Apostle Of Purifying Light")
	card.ManaCost = "{1}{W}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"HUMAN", "CLERIC"}
	card.Power = "2"
	card.Toughness = "1"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddManaCost("{2}").
		// TODO: ExileTargetEffect with complex parameters
		Build()
	card.AddAbility(ability0)
	return card, nil
}
