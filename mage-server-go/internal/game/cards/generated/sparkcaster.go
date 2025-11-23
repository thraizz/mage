package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Sparkcaster", NewSparkcaster)
}

// NewSparkcaster creates a Sparkcaster
// {2}{R}{G} - CREATURE
func NewSparkcaster(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Sparkcaster")
	card.ManaCost = "{2}{R}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"KAVU"}
	card.Power = "5"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewDamageEffect(1)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
