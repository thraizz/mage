package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Jandors Saddlebags", NewJandorsSaddlebags)
}

// NewJandorsSaddlebags creates a Jandors Saddlebags
// {2} - ARTIFACT
func NewJandorsSaddlebags(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Jandors Saddlebags")
	card.ManaCost = "{2}"
	card.Types = []string{"ARTIFACT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddManaCost("{3}").
		AddTapCost().
		// TODO: UntapTargetEffect with complex parameters
		Build()
	card.AddAbility(ability0)
	return card, nil
}
