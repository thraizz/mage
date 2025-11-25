package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Ice Cauldron", NewIceCauldron)
}

// NewIceCauldron creates a Ice Cauldron
// {4} - ARTIFACT
func NewIceCauldron(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Ice Cauldron")
	card.ManaCost = "{4}"
	card.Types = []string{"ARTIFACT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement triggered ability: ActivateIfConditionActivatedAbility
	//   - Effect: IceCauldronExileEffect()
	// card.AddAbility(ability0)
	// TODO: Implement triggered ability: SimpleManaAbility
	//   - Effect: IceCauldronAddManaEffect()
	// card.AddAbility(ability1)
	return card, nil
}
