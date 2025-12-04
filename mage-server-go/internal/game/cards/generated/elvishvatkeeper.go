package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Elvish Vatkeeper", NewElvishVatkeeper)
}

// NewElvishVatkeeper creates a Elvish Vatkeeper
// {1}{B}{G} - CREATURE
func NewElvishVatkeeper(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Elvish Vatkeeper")
	card.ManaCost = "{1}{B}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"PHYREXIAN", "ELF"}
	card.Power = "3"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement activated ability with unmapped effects
	//   - TransformTargetEffect()
	//
	// Costs:
	//   - AddManaCost("{5}")
	// card.AddAbility(ability0)
	return card, nil
}
