package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/counters"
)

func init() {
	cards.Register("Solphim Mayhem Dominus", NewSolphimMayhemDominus)
}

// NewSolphimMayhemDominus creates a Solphim Mayhem Dominus
// {2}{R}{R} - CREATURE
func NewSolphimMayhemDominus(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Solphim Mayhem Dominus")
	card.ManaCost = "{2}{R}{R}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"PHYREXIAN", "HORROR"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "5"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddEffect(abilities.NewAddCountersSourceEffect(counters.NewCounter("indestructible", 1))).
		Build()
	card.AddAbility(ability0)
	return card, nil
}
