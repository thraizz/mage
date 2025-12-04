package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Hua Tuo Honored Physician", NewHuaTuoHonoredPhysician)
}

// NewHuaTuoHonoredPhysician creates a Hua Tuo Honored Physician
// {1}{G}{G} - CREATURE
func NewHuaTuoHonoredPhysician(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Hua Tuo Honored Physician")
	card.ManaCost = "{1}{G}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"HUMAN"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "1"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement triggered ability: ActivateIfConditionActivatedAbility
	//   - Effect: PutOnLibraryTargetEffect()
	// card.AddAbility(ability0)
	return card, nil
}
