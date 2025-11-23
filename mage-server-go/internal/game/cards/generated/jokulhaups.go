package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Jokulhaups", NewJokulhaups)
}

// NewJokulhaups creates a Jokulhaups
// {4}{R}{R} - SORCERY
func NewJokulhaups(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Jokulhaups")
	card.ManaCost = "{4}{R}{R}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewDestroyEffect(filter, true)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
