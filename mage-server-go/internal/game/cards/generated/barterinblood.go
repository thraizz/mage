package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Barter In Blood", NewBarterInBlood)
}

// NewBarterInBlood creates a Barter In Blood
// {2}{B}{B} - SORCERY
func NewBarterInBlood(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Barter In Blood")
	card.ManaCost = "{2}{B}{B}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - SacrificeAllEffect(2, StaticFilters.FILTER_PERMANENT_CREATURES)
	// card.AddAbility(ability0)
	return card, nil
}