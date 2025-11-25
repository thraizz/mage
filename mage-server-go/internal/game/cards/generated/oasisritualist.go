package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Oasis Ritualist", NewOasisRitualist)
}

// NewOasisRitualist creates a Oasis Ritualist
// {3}{G} - CREATURE
func NewOasisRitualist(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Oasis Ritualist")
	card.ManaCost = "{3}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"SNAKE", "DRUID"}
	card.Power = "2"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement triggered ability: SimpleManaAbility
	//   - Effect: AddManaOfAnyColorEffect()
	// card.AddAbility(ability0)
	return card, nil
}
