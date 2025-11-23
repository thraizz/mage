package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Jagged Scar Archers", NewJaggedScarArchers)
}

// NewJaggedScarArchers creates a Jagged Scar Archers
// {1}{G}{G} - CREATURE
func NewJaggedScarArchers(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Jagged Scar Archers")
	card.ManaCost = "{1}{G}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"ELF", "ARCHER"}
	card.Power = "0"
	card.Toughness = "0"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddTapCost().
		AddEffect(abilities.NewDamageEffect(SourcePermanentPowerValue.NOT_NEGATIVE)).
		Build()
	card.AddAbility(ability0)
	return card, nil
}
