package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/token"
)

func init() {
	cards.Register("Minsc Boo Timeless Heroes", NewMinscBooTimelessHeroes)
}

// NewMinscBooTimelessHeroes creates a Minsc Boo Timeless Heroes
// {2}{R}{G} - PLANESWALKER
func NewMinscBooTimelessHeroes(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Minsc Boo Timeless Heroes")
	card.ManaCost = "{2}{R}{G}"
	card.Types = []string{"PLANESWALKER"}
	card.Subtypes = []string{"MINSC"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Loyalty = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	token0_0, err := token.GetToken("BooToken")
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
