package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/effects"
)

func init() {
	cards.Register("Kessig Dire Swine", NewKessigDireSwine)
}

// NewKessigDireSwine creates a Kessig Dire Swine
// {4}{G}{G} - CREATURE
func NewKessigDireSwine(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Kessig Dire Swine")
	card.ManaCost = "{4}{G}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"BOAR", "HORROR"}
	card.Power = "6"
	card.Toughness = "6"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewGrantAbilityEffect("TrampleAbility", effects.DurationWhileOnBattlefield)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
