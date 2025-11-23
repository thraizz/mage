package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Clouds Limit Break", NewCloudsLimitBreak)
}

// NewCloudsLimitBreak creates a Clouds Limit Break
// {1}{W} - INSTANT
func NewCloudsLimitBreak(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Clouds Limit Break")
	card.ManaCost = "{1}{W}"
	card.Types = []string{"INSTANT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewDestroyEffect()).
		AddEffect(abilities.NewDestroyEffect()).
		AddEffect(abilities.NewDestroyEffect(filter)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
