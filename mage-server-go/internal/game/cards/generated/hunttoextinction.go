package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/counters"
)

func init() {
	cards.Register("Hunt To Extinction", NewHuntToExtinction)
}

// NewHuntToExtinction creates a Hunt To Extinction
// {X}{B}{R}{G} - SORCERY
func NewHuntToExtinction(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Hunt To Extinction")
	card.ManaCost = "{X}{B}{R}{G}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - DamageAllEffect(GetXValue.instance, filter)
	//   - DamageAllEffect(GetXValue.instance, new FilterCreaturePermanent())
	// card.AddAbility(ability0)
	return card, nil
}
