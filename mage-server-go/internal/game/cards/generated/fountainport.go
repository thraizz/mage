package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/token"
)

func init() {
	cards.Register("Fountainport", NewFountainport)
}

// NewFountainport creates a Fountainport
//   - LAND
func NewFountainport(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Fountainport")
	card.ManaCost = ""
	card.Types = []string{"LAND"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.BuildSimpleManaAbility(card.ID, "C")
	card.AddAbility(ability0)
	ability1 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddManaCost("{2}").
		AddTapCost().
		AddEffect(abilities.NewDrawCardsEffect(1)).
		Build()
	card.AddAbility(ability1)
	token2_0, err := token.GetToken("FishNoAbilityToken")
	if err != nil {
		return nil, err
	}
	ability2 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddManaCost("{3}").
		AddTapCost().
		AddEffect(abilities.NewCreateTokenEffect(token2_0)).
		Build()
	card.AddAbility(ability2)
	token3_0, err := token.GetToken("TreasureToken")
	if err != nil {
		return nil, err
	}
	ability3 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddManaCost("{4}").
		AddTapCost().
		AddEffect(abilities.NewCreateTokenEffect(token3_0)).
		Build()
	card.AddAbility(ability3)
	return card, nil
}
