package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Interdisciplinary Mascot", NewInterdisciplinaryMascot)
}

// NewInterdisciplinaryMascot creates a Interdisciplinary Mascot
// {6}{U}{U} - CREATURE
func NewInterdisciplinaryMascot(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Interdisciplinary Mascot")
	card.ManaCost = "{6}{U}{U}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"ELEMENTAL", "FRACTAL"}
	card.Power = "5"
	card.Toughness = "5"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - LookLibraryAndPickControllerEffect(                 4, 1, PutCards.HAND, PutCards.BOT...)
	// card.AddAbility(ability0)
	return card, nil
}
