package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/token"
)

func init() {
	cards.Register("Jan Jansen Chaos Crafter", NewJanJansenChaosCrafter)
}

// NewJanJansenChaosCrafter creates a Jan Jansen Chaos Crafter
// {R}{W}{B} - CREATURE
// Haste
func NewJanJansenChaosCrafter(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Jan Jansen Chaos Crafter")
	card.ManaCost = "{R}{W}{B}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"GNOME", "ARTIFICER"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "3"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordHaste)
	card.AddAbility(ability0)
	token1_0, err := token.GetToken("TreasureToken")
	if err != nil {
		return nil, err
	}
	ability1 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddTapCost().
		AddEffect(abilities.NewCreateTokenEffect(token1_0, 2)).
		Build()
	card.AddAbility(ability1)
	token2_0, err := token.GetToken("ConstructToken")
	if err != nil {
		return nil, err
	}
	ability2 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddTapCost().
		AddEffect(abilities.NewCreateTokenEffect(token2_0, 2)).
		Build()
	card.AddAbility(ability2)
	return card, nil
}
