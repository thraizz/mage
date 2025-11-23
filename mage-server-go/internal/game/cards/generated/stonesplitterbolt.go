package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Stonesplitter Bolt", NewStonesplitterBolt)
}

// NewStonesplitterBolt creates a Stonesplitter Bolt
// {X}{R} - INSTANT
func NewStonesplitterBolt(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Stonesplitter Bolt")
	card.ManaCost = "{X}{R}"
	card.Types = []string{"INSTANT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewDamageEffect(new MultipliedValue(GetXValue.instance, 2))).
		AddEffect(abilities.NewDamageEffect(GetXValue.instance)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}