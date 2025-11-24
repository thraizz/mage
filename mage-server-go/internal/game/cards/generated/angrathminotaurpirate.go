package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Angrath Minotaur Pirate", NewAngrathMinotaurPirate)
}

// NewAngrathMinotaurPirate creates a Angrath Minotaur Pirate
// {4}{B}{R} - PLANESWALKER
func NewAngrathMinotaurPirate(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Angrath Minotaur Pirate")
	card.ManaCost = "{4}{B}{R}"
	card.Types = []string{"PLANESWALKER"}
	card.Subtypes = []string{"ANGRATH"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Loyalty = "5"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - ReturnFromGraveyardToBattlefieldTargetEffect()
	// card.AddAbility(ability0)
	return card, nil
}
