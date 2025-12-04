package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("The Capitoline Triad", NewTheCapitolineTriad)
}

// NewTheCapitolineTriad creates a The Capitoline Triad
// {10} - CREATURE
func NewTheCapitolineTriad(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "The Capitoline Triad")
	card.ManaCost = "{10}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"GOD", "ARTIFICER"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "7"
	card.Toughness = "7"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement activated ability with unmapped effects
	//   - GetEmblemEffect()
	// card.AddAbility(ability0)
	return card, nil
}
