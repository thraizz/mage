package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Subversive Acolyte", NewSubversiveAcolyte)
}

// NewSubversiveAcolyte creates a Subversive Acolyte
// {1}{B} - CREATURE
func NewSubversiveAcolyte(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Subversive Acolyte")
	card.ManaCost = "{1}{B}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"HUMAN"}
	card.Power = "2"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement triggered ability: ActivateOncePerGameActivatedAbility
	//   - Effect: AddCardSubTypeSourceEffect()
	// card.AddAbility(ability0)
	return card, nil
}
