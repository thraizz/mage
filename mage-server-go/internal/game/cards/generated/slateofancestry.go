package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Slate Of Ancestry", NewSlateOfAncestry)
}

// NewSlateOfAncestry creates a Slate Of Ancestry
// {4} - ARTIFACT
func NewSlateOfAncestry(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Slate Of Ancestry")
	card.ManaCost = "{4}"
	card.Types = []string{"ARTIFACT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddManaCost("{4}").
		AddTapCost().
		// TODO: DrawCardSourceControllerEffect with complex parameters
		Build()
	card.AddAbility(ability0)
	return card, nil
}
