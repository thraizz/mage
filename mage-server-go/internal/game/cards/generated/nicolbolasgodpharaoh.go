package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Nicol Bolas God Pharaoh", NewNicolBolasGodPharaoh)
}

// NewNicolBolasGodPharaoh creates a Nicol Bolas God Pharaoh
// {4}{U}{B}{R} - PLANESWALKER
func NewNicolBolasGodPharaoh(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Nicol Bolas God Pharaoh")
	card.ManaCost = "{4}{U}{B}{R}"
	card.Types = []string{"PLANESWALKER"}
	card.Subtypes = []string{"BOLAS"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Loyalty = "7"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewDamageEffect(7)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	ability1, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		// TODO: ExileAllEffect with complex parameters
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability1)
	return card, nil
}
