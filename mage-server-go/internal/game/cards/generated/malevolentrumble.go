package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/token"
)

func init() {
	cards.Register("Malevolent Rumble", NewMalevolentRumble)
}

// NewMalevolentRumble creates a Malevolent Rumble
// {1}{G} - SORCERY
func NewMalevolentRumble(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Malevolent Rumble")
	card.ManaCost = "{1}{G}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - RevealLibraryPickControllerEffect(                 4, 1, StaticFilters.FILTER_CARD_P...)
	// card.AddAbility(ability0)
	return card, nil
}
