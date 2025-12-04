package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Corrosive Gale", NewCorrosiveGale)
}

// NewCorrosiveGale creates a Corrosive Gale
// {X}{G/P} - SORCERY
func NewCorrosiveGale(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Corrosive Gale")
	card.ManaCost = "{X}{G/P}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - DamageAllEffect(GetXValue.instance, StaticFilters.FILTER_CREATURE_...)
	// card.AddAbility(ability0)
	return card, nil
}
