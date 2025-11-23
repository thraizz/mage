package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/token"
)

func init() {
	cards.Register("Izoni Center Of The Web", NewIzoniCenterOfTheWeb)
}

// NewIzoniCenterOfTheWeb creates a Izoni Center Of The Web
// {4}{B}{G} - CREATURE
func NewIzoniCenterOfTheWeb(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Izoni Center Of The Web")
	card.ManaCost = "{4}{B}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"ELF", "DETECTIVE"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "5"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddEffect(abilities.NewSurveilEffect(1)).
		AddEffect(abilities.NewDrawCardsEffect(2)).
		AddEffect(abilities.NewGainLifeEffect(2)).
		Build()
	card.AddAbility(ability0)
	// TODO: Implement spell ability with unmapped effects
	//   - DoIfCostPaid(                 new CreateTokenEffect(new IzoniSp...)
	// card.AddAbility(ability1)
	return card, nil
}
