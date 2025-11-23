package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Ravenous Demon", NewRavenousDemon)
}

// NewRavenousDemon creates a Ravenous Demon
// {3}{B}{B} - CREATURE
func NewRavenousDemon(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Ravenous Demon")
	card.ManaCost = "{3}{B}{B}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"DEMON"}
	card.Power = "4"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - TransformSourceEffect()
	// card.AddAbility(ability0)
	return card, nil
}
