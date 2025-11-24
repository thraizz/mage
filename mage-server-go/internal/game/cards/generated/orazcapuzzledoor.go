package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Orazca Puzzle Door", NewOrazcaPuzzleDoor)
}

// NewOrazcaPuzzleDoor creates a Orazca Puzzle Door
// {U} - ARTIFACT
func NewOrazcaPuzzleDoor(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Orazca Puzzle Door")
	card.ManaCost = "{U}"
	card.Types = []string{"ARTIFACT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement activated ability with unmapped effects
	//   - LookLibraryAndPickControllerEffect(                 2, 1, PutCards.HAND, PutCards.GRA...)
	//
	// Costs:
	//   - AddManaCost("{1}")
	//   - AddTapCost()
	//   - AddSacrificeSourceCost()
	// card.AddAbility(ability0)
	return card, nil
}
