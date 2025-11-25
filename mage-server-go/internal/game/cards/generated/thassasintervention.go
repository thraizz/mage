package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/counters"
)

func init() {
	cards.Register("Thassas Intervention", NewThassasIntervention)
}

// NewThassasIntervention creates a Thassas Intervention
// {X}{U}{U} - INSTANT
func NewThassasIntervention(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Thassas Intervention")
	card.ManaCost = "{X}{U}{U}"
	card.Types = []string{"INSTANT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - LookLibraryAndPickControllerEffect(                 GetXValue.instance, 2, PutCards.H...)
	// card.AddAbility(ability0)
	return card, nil
}
