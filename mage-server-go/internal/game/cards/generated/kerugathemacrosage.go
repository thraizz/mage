package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Keruga The Macrosage", NewKerugaTheMacrosage)
}

// NewKerugaTheMacrosage creates a Keruga The Macrosage
// {3}{G/U}{G/U} - CREATURE
func NewKerugaTheMacrosage(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Keruga The Macrosage")
	card.ManaCost = "{3}{G/U}{G/U}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"DINOSAUR", "HIPPO"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "5"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewDrawCardsEffect()).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
