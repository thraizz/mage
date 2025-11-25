package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Ellie And Alan Paleontologists", NewEllieAndAlanPaleontologists)
}

// NewEllieAndAlanPaleontologists creates a Ellie And Alan Paleontologists
// {2}{G}{W}{U} - CREATURE
func NewEllieAndAlanPaleontologists(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Ellie And Alan Paleontologists")
	card.ManaCost = "{2}{G}{W}{U}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"HUMAN", "SCIENTIST"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "2"
	card.Toughness = "5"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement triggered ability: ActivateAsSorceryActivatedAbility
	//   - Effect: EllieAndAlanDiscoverEffect()
	// card.AddAbility(ability0)
	return card, nil
}
