package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Kozileks Command", NewKozileksCommand)
}

// NewKozileksCommand creates a Kozileks Command
// {X}{C}{C} - KINDRED INSTANT
func NewKozileksCommand(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Kozileks Command")
	card.ManaCost = "{X}{C}{C}"
	card.Types = []string{"KINDRED", "INSTANT"}
	card.Subtypes = []string{"ELDRAZI"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		// TODO: ExileTargetEffect with complex parameters
		// TODO: ExileTargetEffect with complex parameters
		AddTarget(abilities.NewPlayerTargetFilter()).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
