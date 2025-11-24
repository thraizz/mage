package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/token"
)

func init() {
	cards.Register("Eyes Of The Wisent", NewEyesOfTheWisent)
}

// NewEyesOfTheWisent creates a Eyes Of The Wisent
// {1}{G} - KINDRED ENCHANTMENT
func NewEyesOfTheWisent(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Eyes Of The Wisent")
	card.ManaCost = "{1}{G}"
	card.Types = []string{"KINDRED", "ENCHANTMENT"}
	card.Subtypes = []string{"ELEMENTAL"}
	card.SetCode = "M21"
	card.Rarity = "common"

	token0_0, err := token.GetToken("Elemental44GreenToken")
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
