package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Mask Of The Jadecrafter", NewMaskOfTheJadecrafter)
}

// NewMaskOfTheJadecrafter creates a Mask Of The Jadecrafter
// {2} - ARTIFACT
func NewMaskOfTheJadecrafter(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Mask Of The Jadecrafter")
	card.ManaCost = "{2}"
	card.Types = []string{"ARTIFACT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement triggered ability: ActivateAsSorceryActivatedAbility
	//   - Effect: MaskOfTheJadecrafterEffect()
	// card.AddAbility(ability0)
	return card, nil
}
