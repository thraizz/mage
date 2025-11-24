package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Argivian Archaeologist", NewArgivianArchaeologist)
}

// NewArgivianArchaeologist creates a Argivian Archaeologist
// {1}{W}{W} - CREATURE
func NewArgivianArchaeologist(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Argivian Archaeologist")
	card.ManaCost = "{1}{W}{W}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"HUMAN", "ARTIFICER"}
	card.Power = "1"
	card.Toughness = "1"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddTapCost().
		// TODO: ReturnFromGraveyardToHandTargetEffect with complex parameters
		Build()
	card.AddAbility(ability0)
	return card, nil
}
