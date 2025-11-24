package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Wine Of Blood And Iron", NewWineOfBloodAndIron)
}

// NewWineOfBloodAndIron creates a Wine Of Blood And Iron
// {3} - ARTIFACT
func NewWineOfBloodAndIron(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Wine Of Blood And Iron")
	card.ManaCost = "{3}"
	card.Types = []string{"ARTIFACT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddManaCost("{4}").
		AddEffect(abilities.NewBoostEffect(TargetPermanentPowerCount.instance, StaticValue.get(0))).
		Build()
	card.AddAbility(ability0)
	return card, nil
}