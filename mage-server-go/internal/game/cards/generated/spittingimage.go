package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/token"
)

func init() {
	cards.Register("Spitting Image", NewSpittingImage)
}

// NewSpittingImage creates a Spitting Image
// {4}{G/U}{G/U} - SORCERY
func NewSpittingImage(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Spitting Image")
	card.ManaCost = "{4}{G/U}{G/U}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - CreateTokenCopyTargetEffect()
	//
	// Targets:
	//   - abilities.NewCreatureTargetFilter()
	// card.AddAbility(ability0)
	return card, nil
}
