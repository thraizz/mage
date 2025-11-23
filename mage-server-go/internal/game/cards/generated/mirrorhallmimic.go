package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Mirrorhall Mimic", NewMirrorhallMimic)
}

// NewMirrorhallMimic creates a Mirrorhall Mimic
// {3}{U} - CREATURE
func NewMirrorhallMimic(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Mirrorhall Mimic")
	card.ManaCost = "{3}{U}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"SPIRIT"}
	card.Power = "0"
	card.Toughness = "0"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - CopyPermanentEffect(                 StaticFilters.FILTER_PERMANENT_CR...)
	// card.AddAbility(ability0)
	return card, nil
}
