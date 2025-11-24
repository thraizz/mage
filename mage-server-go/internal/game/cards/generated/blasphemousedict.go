package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Blasphemous Edict", NewBlasphemousEdict)
}

// NewBlasphemousEdict creates a Blasphemous Edict
// {3}{B}{B} - SORCERY
func NewBlasphemousEdict(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Blasphemous Edict")
	card.ManaCost = "{3}{B}{B}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - SacrificeAllEffect(13, StaticFilters.FILTER_PERMANENT_CREATURES)
	// card.AddAbility(ability0)
	return card, nil
}