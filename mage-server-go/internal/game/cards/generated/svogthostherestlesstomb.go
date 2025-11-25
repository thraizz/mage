package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Svogthos The Restless Tomb", NewSvogthosTheRestlessTomb)
}

// NewSvogthosTheRestlessTomb creates a Svogthos The Restless Tomb
//   - LAND
func NewSvogthosTheRestlessTomb(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Svogthos The Restless Tomb")
	card.ManaCost = ""
	card.Types = []string{"LAND"}
	card.Subtypes = []string{"PLANT", "ZOMBIE"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.BuildSimpleManaAbility(card.ID, "C")
	card.AddAbility(ability0)
	// TODO: Implement activated ability with unmapped effects
	//   - BecomesCreatureSourceEffect()
	// card.AddAbility(ability1)
	return card, nil
}
