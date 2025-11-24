package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Sakashima The Impostor", NewSakashimaTheImpostor)
}

// NewSakashimaTheImpostor creates a Sakashima The Impostor
// {2}{U}{U} - CREATURE
func NewSakashimaTheImpostor(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Sakashima The Impostor")
	card.ManaCost = "{2}{U}{U}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"HUMAN", "ROGUE"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "3"
	card.Toughness = "1"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - CopyPermanentEffect(StaticFilters.FILTER_PERMANENT_CREATURE, new Sakas...)
	// card.AddAbility(ability0)
	return card, nil
}
