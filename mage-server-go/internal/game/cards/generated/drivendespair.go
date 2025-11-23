package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Driven Despair", NewDrivenDespair)
}

// NewDrivenDespair creates a Driven Despair
// {1}{G} - SORCERY
func NewDrivenDespair(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Driven Despair")
	card.ManaCost = "{1}{G}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - DiscardTargetEffect(1)
	// card.AddAbility(ability0)
	return card, nil
}
