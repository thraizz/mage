package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/token"
)

func init() {
	cards.Register("Tevesh Szat Doom Of Fools", NewTeveshSzatDoomOfFools)
}

// NewTeveshSzatDoomOfFools creates a Tevesh Szat Doom Of Fools
// {4}{B} - PLANESWALKER
func NewTeveshSzatDoomOfFools(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Tevesh Szat Doom Of Fools")
	card.ManaCost = "{4}{B}"
	card.Types = []string{"PLANESWALKER"}
	card.Subtypes = []string{"SZAT"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Loyalty = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	token0_0, err := token.GetToken("BreedingPitThrullToken")
	if err != nil {
		return nil, err
	}
	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewCreateTokenEffect(token0_0, 2)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
