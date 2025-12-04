package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Portal Of Sanctuary", NewPortalOfSanctuary)
}

// NewPortalOfSanctuary creates a Portal Of Sanctuary
// {2}{U} - ARTIFACT
func NewPortalOfSanctuary(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Portal Of Sanctuary")
	card.ManaCost = "{2}{U}"
	card.Types = []string{"ARTIFACT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement triggered ability: ActivateIfConditionActivatedAbility
	//   - Effect: PortalOfSanctuaryEffect()
	// card.AddAbility(ability0)
	return card, nil
}
