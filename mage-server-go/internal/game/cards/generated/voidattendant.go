package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/token"
)

func init() {
	cards.Register("Void Attendant", NewVoidAttendant)
}

// NewVoidAttendant creates a Void Attendant
// {2}{G} - CREATURE
func NewVoidAttendant(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Void Attendant")
	card.ManaCost = "{2}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"ELDRAZI", "PROCESSOR"}
	card.Power = "2"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	token0_0, err := token.GetToken("EldraziScionToken")
	if err != nil {
		return nil, err
	}
	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewCreateTokenEffect(token0_0)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
