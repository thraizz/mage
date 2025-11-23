package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/counters"
)

func init() {
	cards.Register("Ayesha Tanaka", NewAyeshaTanaka)
}

// NewAyeshaTanaka creates a Ayesha Tanaka
// {W}{W}{U}{U} - CREATURE
func NewAyeshaTanaka(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Ayesha Tanaka")
	card.ManaCost = "{W}{W}{U}{U}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"HUMAN", "ARTIFICER"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "2"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddTapCost().
		AddEffect(abilities.NewCounterSpellEffect(new ManaCostsImpl<>("{W}"))).
		Build()
	card.AddAbility(ability0)
	return card, nil
}