package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Brilliant Spectrum", NewBrilliantSpectrum)
}

// NewBrilliantSpectrum creates a Brilliant Spectrum
// {3}{U} - SORCERY
func NewBrilliantSpectrum(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Brilliant Spectrum")
	card.ManaCost = "{3}{U}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - DiscardControllerEffect(2)
	// card.AddAbility(ability0)
	return card, nil
}
