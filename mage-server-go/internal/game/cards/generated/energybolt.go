package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Energy Bolt", NewEnergyBolt)
}

// NewEnergyBolt creates a Energy Bolt
// {X}{R}{W} - SORCERY
func NewEnergyBolt(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Energy Bolt")
	card.ManaCost = "{X}{R}{W}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewGainLifeEffect(GetXValue.instance)).
		AddEffect(abilities.NewDamageEffect(GetXValue.instance)).
		AddEffect(abilities.NewGainLifeEffect(GetXValue.instance)).
		AddTarget(abilities.NewTargetRequirement(1, 1, abilities.NewAnyTargetFilter())).
		AddTarget(abilities.NewTargetRequirement(1, 1, abilities.NewPlayerTargetFilter())).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
