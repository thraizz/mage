package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Masterwork Of Ingenuity", NewMasterworkOfIngenuity)
}

// NewMasterworkOfIngenuity creates a Masterwork Of Ingenuity
// {1} - ARTIFACT
func NewMasterworkOfIngenuity(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Masterwork Of Ingenuity")
	card.ManaCost = "{1}"
	card.Types = []string{"ARTIFACT"}
	card.Subtypes = []string{"EQUIPMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - CopyPermanentEffect(StaticFilters.FILTER_PERMANENT_EQUIPMENT)
	// card.AddAbility(ability0)
	return card, nil
}
