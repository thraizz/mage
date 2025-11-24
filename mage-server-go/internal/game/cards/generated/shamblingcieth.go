package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Shambling Cieth", NewShamblingCieth)
}

// NewShamblingCieth creates a Shambling Cieth
// {2}{B} - CREATURE
func NewShamblingCieth(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Shambling Cieth")
	card.ManaCost = "{2}{B}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"MUTANT", "HORROR"}
	card.Power = "3"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - DoIfCostPaid(new ReturnSourceFromGraveyardToHandEffect(), new M...)
	// card.AddAbility(ability0)
	return card, nil
}
