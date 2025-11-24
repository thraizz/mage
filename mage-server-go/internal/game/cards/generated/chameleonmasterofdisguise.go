package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Chameleon Master Of Disguise", NewChameleonMasterOfDisguise)
}

// NewChameleonMasterOfDisguise creates a Chameleon Master Of Disguise
// {3}{U} - CREATURE
func NewChameleonMasterOfDisguise(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Chameleon Master Of Disguise")
	card.ManaCost = "{3}{U}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"HUMAN", "SHAPESHIFTER", "VILLAIN"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "2"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - CopyPermanentEffect(                 StaticFilters.FILTER_CONTROLLED_C...)
	// card.AddAbility(ability0)
	return card, nil
}
