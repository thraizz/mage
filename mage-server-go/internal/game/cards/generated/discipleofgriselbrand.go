package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Disciple Of Griselbrand", NewDiscipleOfGriselbrand)
}

// NewDiscipleOfGriselbrand creates a Disciple Of Griselbrand
// {1}{B} - CREATURE
func NewDiscipleOfGriselbrand(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Disciple Of Griselbrand")
	card.ManaCost = "{1}{B}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"HUMAN", "CLERIC"}
	card.Power = "1"
	card.Toughness = "1"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement activated ability with unmapped effects
	//   - DiscipleOfGriselbrandEffect()
	//
	// Costs:
	//   - AddManaCost("{1}")
	// card.AddAbility(ability0)
	return card, nil
}
