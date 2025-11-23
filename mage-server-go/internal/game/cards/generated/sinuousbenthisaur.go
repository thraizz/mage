package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Sinuous Benthisaur", NewSinuousBenthisaur)
}

// NewSinuousBenthisaur creates a Sinuous Benthisaur
// {5}{U} - CREATURE
func NewSinuousBenthisaur(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Sinuous Benthisaur")
	card.ManaCost = "{5}{U}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"DINOSAUR"}
	card.Power = "4"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - LookLibraryAndPickControllerEffect(                         CavesControlledAndInGrave...)
	// card.AddAbility(ability0)
	return card, nil
}
