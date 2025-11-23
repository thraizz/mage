package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Simulacrum", NewSimulacrum)
}

// NewSimulacrum creates a Simulacrum
// {1}{B} - INSTANT
func NewSimulacrum(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Simulacrum")
	card.ManaCost = "{1}{B}"
	card.Types = []string{"INSTANT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewDamageEffect(new SimulacrumAmount(), true, true)).
		AddEffect(abilities.NewGainLifeEffect(new SimulacrumAmount())).
		AddEffect(abilities.NewDamageEffect(new SimulacrumAmount(), true, true)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}