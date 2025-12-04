package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Genesis Ultimatum", NewGenesisUltimatum)
}

// NewGenesisUltimatum creates a Genesis Ultimatum
// {G}{G}{U}{U}{U}{R}{R} - SORCERY
func NewGenesisUltimatum(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Genesis Ultimatum")
	card.ManaCost = "{G}{G}{U}{U}{U}{R}{R}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - LookLibraryAndPickControllerEffect(                 5, Integer.MAX_VALUE, filter, Put...)
	// card.AddAbility(ability0)
	return card, nil
}
