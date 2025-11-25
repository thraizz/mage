package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Adventurous Impulse", NewAdventurousImpulse)
}

// NewAdventurousImpulse creates a Adventurous Impulse
// {G} - SORCERY
func NewAdventurousImpulse(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Adventurous Impulse")
	card.ManaCost = "{G}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - LookLibraryAndPickControllerEffect(                 3, 1, StaticFilters.FILTER_CARD_C...)
	// card.AddAbility(ability0)
	return card, nil
}
