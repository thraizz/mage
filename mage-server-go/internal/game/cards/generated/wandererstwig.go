package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Wanderers Twig", NewWanderersTwig)
}

// NewWanderersTwig creates a Wanderers Twig
// {1} - ARTIFACT
func NewWanderersTwig(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Wanderers Twig")
	card.ManaCost = "{1}"
	card.Types = []string{"ARTIFACT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddManaCost("{1}").
		AddSacrificeSourceCost().
		// TODO: SearchLibraryPutInHandEffect with complex parameters
		Build()
	card.AddAbility(ability0)
	return card, nil
}
