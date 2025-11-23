package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Kitsune Diviner", NewKitsuneDiviner)
}

// NewKitsuneDiviner creates a Kitsune Diviner
// {W} - CREATURE
func NewKitsuneDiviner(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Kitsune Diviner")
	card.ManaCost = "{W}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"FOX", "CLERIC"}
	card.Power = "0"
	card.Toughness = "1"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddTapCost().
		AddEffect(abilities.NewTapEffect()).
		Build()
	card.AddAbility(ability0)
	return card, nil
}
