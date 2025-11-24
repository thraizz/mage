package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Necrotic Hex", NewNecroticHex)
}

// NewNecroticHex creates a Necrotic Hex
// {6}{B} - SORCERY
func NewNecroticHex(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Necrotic Hex")
	card.ManaCost = "{6}{B}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - SacrificeAllEffect(6, StaticFilters.FILTER_PERMANENT_CREATURES)
	// card.AddAbility(ability0)
	return card, nil
}