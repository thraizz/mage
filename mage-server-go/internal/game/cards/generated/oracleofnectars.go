package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Oracle Of Nectars", NewOracleOfNectars)
}

// NewOracleOfNectars creates a Oracle Of Nectars
// {2}{G/W} - CREATURE
func NewOracleOfNectars(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Oracle Of Nectars")
	card.ManaCost = "{2}{G/W}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"ELF", "CLERIC"}
	card.Power = "2"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddTapCost().
		// TODO: GainLifeEffect with complex parameters
		Build()
	card.AddAbility(ability0)
	return card, nil
}
