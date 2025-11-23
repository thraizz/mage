package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Tidal Wave", NewTidalWave)
}

// NewTidalWave creates a Tidal Wave
// {2}{U} - INSTANT
func NewTidalWave(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Tidal Wave")
	card.ManaCost = "{2}{U}"
	card.Types = []string{"INSTANT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - SacrificeTargetEffect()
	// card.AddAbility(ability0)
	return card, nil
}
