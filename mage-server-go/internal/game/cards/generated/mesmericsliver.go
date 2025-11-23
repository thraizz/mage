package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/effects"
)

func init() {
	cards.Register("Mesmeric Sliver", NewMesmericSliver)
}

// NewMesmericSliver creates a Mesmeric Sliver
// {3}{U} - CREATURE
func NewMesmericSliver(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Mesmeric Sliver")
	card.ManaCost = "{3}{U}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"SLIVER"}
	card.Power = "2"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewGrantAbilityEffect(filterSliverWhen permanent enters the battlefield, you may fateseal 1.\" <i>(To fateseal 1, its controller looks at the top card of an opponent's library, then they may put that card on the bottom of that library.)</i>")).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}