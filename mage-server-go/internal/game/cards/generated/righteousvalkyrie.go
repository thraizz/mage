package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Righteous Valkyrie", NewRighteousValkyrie)
}

// NewRighteousValkyrie creates a Righteous Valkyrie
// {2}{W} - CREATURE
// Flying
func NewRighteousValkyrie(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Righteous Valkyrie")
	card.ManaCost = "{2}{W}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"ANGEL", "CLERIC"}
	card.Power = "2"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	ability1, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewBoostEffect(2, 2)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability1)
	return card, nil
}
