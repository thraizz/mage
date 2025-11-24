package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/token"
)

func init() {
	cards.Register("Ashiok Nightmare Muse", NewAshiokNightmareMuse)
}

// NewAshiokNightmareMuse creates a Ashiok Nightmare Muse
// {3}{U}{B} - PLANESWALKER
func NewAshiokNightmareMuse(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Ashiok Nightmare Muse")
	card.ManaCost = "{3}{U}{B}"
	card.Types = []string{"PLANESWALKER"}
	card.Subtypes = []string{"ASHIOK"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Loyalty = "5"
	card.SetCode = "M21"
	card.Rarity = "common"

	token0_0, err := token.GetToken("AshiokNightmareMuseToken")
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