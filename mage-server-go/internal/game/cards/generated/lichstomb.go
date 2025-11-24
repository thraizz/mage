package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Lichs Tomb", NewLichsTomb)
}

// NewLichsTomb creates a Lichs Tomb
// {4} - ARTIFACT
func NewLichsTomb(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Lichs Tomb")
	card.ManaCost = "{4}"
	card.Types = []string{"ARTIFACT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - SacrificeControllerEffect(                 StaticFilters.FILTER_PERMANENT, S...)
	// card.AddAbility(ability0)
	return card, nil
}
