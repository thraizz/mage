package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Barrin Master Wizard", NewBarrinMasterWizard)
}

// NewBarrinMasterWizard creates a Barrin Master Wizard
// {1}{U}{U} - CREATURE
func NewBarrinMasterWizard(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Barrin Master Wizard")
	card.ManaCost = "{1}{U}{U}"
	card.Types = []string{"CREATURE"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "1"
	card.Toughness = "1"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewActivatedAbilityBuilder(card.ID).
		// TODO: ReturnToHandTargetEffect with complex parameters
		Build()
	card.AddAbility(ability0)
	return card, nil
}
