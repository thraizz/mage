package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/effects"
)

func init() {
	cards.Register("Turtle Duck", NewTurtleDuck)
}

// NewTurtleDuck creates a Turtle Duck
// {G} - CREATURE
func NewTurtleDuck(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Turtle Duck")
	card.ManaCost = "{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"TURTLE", "BIRD"}
	card.Power = "0"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddManaCost("{3}").
		AddEffect(abilities.NewGrantAbilityEffect("TrampleAbility", effects.DurationEndOfTurn)).
		Build()
	card.AddAbility(ability0)
	return card, nil
}
