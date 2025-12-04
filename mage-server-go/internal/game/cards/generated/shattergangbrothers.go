package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Shattergang Brothers", NewShattergangBrothers)
}

// NewShattergangBrothers creates a Shattergang Brothers
// {1}{B}{R}{G} - CREATURE
func NewShattergangBrothers(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Shattergang Brothers")
	card.ManaCost = "{1}{B}{R}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"GOBLIN", "ARTIFICER"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "3"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement activated ability with unmapped effects
	//   - ShattergangBrothersEffect()
	// card.AddAbility(ability0)
	// TODO: Implement activated ability with unmapped effects
	//   - ShattergangBrothersEffect()
	// card.AddAbility(ability1)
	// TODO: Implement activated ability with unmapped effects
	//   - ShattergangBrothersEffect()
	// card.AddAbility(ability2)
	return card, nil
}
