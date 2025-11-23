package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/token"
)

func init() {
	cards.Register("Hidden Stockpile", NewHiddenStockpile)
}

// NewHiddenStockpile creates a Hidden Stockpile
// {W}{B} - ENCHANTMENT
func NewHiddenStockpile(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Hidden Stockpile")
	card.ManaCost = "{W}{B}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddManaCost("{1}").
		AddEffect(abilities.NewScryEffect(1)).
		Build()
	card.AddAbility(ability0)
	token1_0, err := token.GetToken("ServoToken")
	if err != nil {
		return nil, err
	}
	ability1, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewCreateTokenEffect(token1_0)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability1)
	return card, nil
}
