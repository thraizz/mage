package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/effects"
)

func init() {
	cards.Register("Moment Of Valor", NewMomentOfValor)
}

// NewMomentOfValor creates a Moment Of Valor
// {2}{W} - INSTANT
func NewMomentOfValor(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Moment Of Valor")
	card.ManaCost = "{2}{W}"
	card.Types = []string{"INSTANT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewDestroyEffect()).
		AddEffect(abilities.NewUntapEffect()).
		AddEffect(abilities.NewBoostEffect(1, 0)).
		AddEffect(abilities.NewGrantAbilityEffect("IndestructibleAbility")).
		AddTarget(abilities.NewCreatureTargetFilter()).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
