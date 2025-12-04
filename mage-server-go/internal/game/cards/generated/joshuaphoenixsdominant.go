package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Joshua Phoenixs Dominant", NewJoshuaPhoenixsDominant)
}

// NewJoshuaPhoenixsDominant creates a Joshua Phoenixs Dominant
// {1}{R}{W} - CREATURE
func NewJoshuaPhoenixsDominant(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Joshua Phoenixs Dominant")
	card.ManaCost = "{1}{R}{W}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"HUMAN", "NOBLE", "WIZARD"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "3"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement triggered ability: ActivateAsSorceryActivatedAbility
	//   - Effect: ExileAndReturnSourceEffect()
	// card.AddAbility(ability0)
	return card, nil
}
