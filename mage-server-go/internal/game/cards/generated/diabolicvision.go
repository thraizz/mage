package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Diabolic Vision", NewDiabolicVision)
}

// NewDiabolicVision creates a Diabolic Vision
// {U}{B} - SORCERY
func NewDiabolicVision(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Diabolic Vision")
	card.ManaCost = "{U}{B}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - LookLibraryAndPickControllerEffect(5, 1, PutCards.HAND, PutCards.TOP_ANY)
	// card.AddAbility(ability0)
	return card, nil
}
