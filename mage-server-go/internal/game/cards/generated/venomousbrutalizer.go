package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Venomous Brutalizer", NewVenomousBrutalizer)
}

// NewVenomousBrutalizer creates a Venomous Brutalizer
// {2}{G}{G} - CREATURE
func NewVenomousBrutalizer(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Venomous Brutalizer")
	card.ManaCost = "{2}{G}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"PHYREXIAN", "KNIGHT"}
	card.Power = "4"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - DoIfCostPaid(new ProliferateEffect(), new ManaCostsImpl<>("{1}{...)
	// card.AddAbility(ability0)
	return card, nil
}
