package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Warden Of The First Tree", NewWardenOfTheFirstTree)
}

// NewWardenOfTheFirstTree creates a Warden Of The First Tree
// {G} - CREATURE
func NewWardenOfTheFirstTree(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Warden Of The First Tree")
	card.ManaCost = "{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"HUMAN"}
	card.Power = "1"
	card.Toughness = "1"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement activated ability with unmapped effects
	//   - AddCardSubTypeSourceEffect()
	// card.AddAbility(ability0)
	return card, nil
}
