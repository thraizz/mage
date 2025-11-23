package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/effects"
)

func init() {
	cards.Register("Plaxcaster Frogling", NewPlaxcasterFrogling)
}

// NewPlaxcasterFrogling creates a Plaxcaster Frogling
// {1}{G}{U} - CREATURE
func NewPlaxcasterFrogling(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Plaxcaster Frogling")
	card.ManaCost = "{1}{G}{U}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"FROG", "MUTANT"}
	card.Power = "0"
	card.Toughness = "0"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddManaCost("{2}").
		AddEffect(abilities.NewGrantAbilityEffect("ShroudAbility", effects.DurationEndOfTurn)).
		Build()
	card.AddAbility(ability0)
	return card, nil
}
