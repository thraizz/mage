package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/counters"
)

func init() {
	cards.Register("Cryptic Trilobite", NewCrypticTrilobite)
}

// NewCrypticTrilobite creates a Cryptic Trilobite
// {X}{X} - CREATURE
func NewCrypticTrilobite(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Cryptic Trilobite")
	card.ManaCost = "{X}{X}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"TRILOBITE"}
	card.Power = "0"
	card.Toughness = "0"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddManaCost("{1}").
		AddTapCost().
		AddEffect(abilities.NewAddCountersSourceEffect(counters.CounterTypeP1P1.CreateInstance(1))).
		Build()
	card.AddAbility(ability0)
	return card, nil
}
