package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/effects"
)

func init() {
	cards.Register("Tyrant Guard", NewTyrantGuard)
}

// NewTyrantGuard creates a Tyrant Guard
// {X}{2}{G} - CREATURE
func NewTyrantGuard(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Tyrant Guard")
	card.ManaCost = "{X}{2}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"TYRANID"}
	card.Power = "3"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddSacrificeSourceCost().
		AddEffect(abilities.NewGrantAbilityEffect("HexproofAbility", effects.DurationEndOfTurn)).
		AddEffect(abilities.NewGrantAbilityEffect("IndestructibleAbility", effects.DurationEndOfTurn)).
		Build()
	card.AddAbility(ability0)
	return card, nil
}
