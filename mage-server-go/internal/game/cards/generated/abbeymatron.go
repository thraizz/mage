package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Abbey Matron", NewAbbeyMatron)
}

// NewAbbeyMatron creates a Abbey Matron
// {2}{W} - CREATURE
func NewAbbeyMatron(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Abbey Matron")
	card.ManaCost = "{2}{W}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"HUMAN", "CLERIC"}
	card.Power = "1"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddTapCost().
		AddEffect(abilities.NewBoostEffect(0, 3)).
		Build()
	card.AddAbility(ability0)
	return card, nil
}
