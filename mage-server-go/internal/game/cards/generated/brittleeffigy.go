package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Brittle Effigy", NewBrittleEffigy)
}

// NewBrittleEffigy creates a Brittle Effigy
// {1} - ARTIFACT
func NewBrittleEffigy(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Brittle Effigy")
	card.ManaCost = "{1}"
	card.Types = []string{"ARTIFACT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddManaCost("{4}").
		AddTapCost().
		AddEffect(abilities.NewExileTargetEffect()).
		AddTarget(abilities.NewCreatureTargetFilter()).
		Build()
	card.AddAbility(ability0)
	return card, nil
}
