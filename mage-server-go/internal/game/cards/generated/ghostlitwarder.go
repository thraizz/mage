package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/counters"
)

func init() {
	cards.Register("Ghost Lit Warder", NewGhostLitWarder)
}

// NewGhostLitWarder creates a Ghost Lit Warder
// {1}{U} - CREATURE
func NewGhostLitWarder(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Ghost Lit Warder")
	card.ManaCost = "{1}{U}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"SPIRIT"}
	card.Power = "1"
	card.Toughness = "1"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewCounterSpellEffect(new ManaCostsImpl<>("{4}"))).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	ability1 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddTapCost().
		AddEffect(abilities.NewCounterSpellEffect(new ManaCostsImpl<>("{2}"))).
		Build()
	card.AddAbility(ability1)
	return card, nil
}