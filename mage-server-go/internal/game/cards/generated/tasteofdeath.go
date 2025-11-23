package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/token"
)

func init() {
	cards.Register("Taste Of Death", NewTasteOfDeath)
}

// NewTasteOfDeath creates a Taste Of Death
// {4}{B}{B} - SORCERY
func NewTasteOfDeath(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Taste Of Death")
	card.ManaCost = "{4}{B}{B}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - SacrificeAllEffect(3, StaticFilters.FILTER_PERMANENT_CREATURES)
	// card.AddAbility(ability0)
	return card, nil
}
