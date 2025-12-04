package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Altar Of Dementia", NewAltarOfDementia)
}

// NewAltarOfDementia creates a Altar Of Dementia
// {2} - ARTIFACT
func NewAltarOfDementia(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Altar Of Dementia")
	card.ManaCost = "{2}"
	card.Types = []string{"ARTIFACT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement activated ability with unmapped effects
	//   - AltarOfDementiaEffect()
	// card.AddAbility(ability0)
	return card, nil
}
