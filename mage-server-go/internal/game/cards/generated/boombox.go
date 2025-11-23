package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Boom Box", NewBoomBox)
}

// NewBoomBox creates a Boom Box
// {2} - ARTIFACT
func NewBoomBox(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Boom Box")
	card.ManaCost = "{2}"
	card.Types = []string{"ARTIFACT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddManaCost("{6}").
		AddTapCost().
		AddSacrificeSourceCost().
		AddEffect(abilities.NewDestroyEffect("")).
		Build()
	card.AddAbility(ability0)
	return card, nil
}
