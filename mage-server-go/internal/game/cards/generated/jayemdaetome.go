package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Jayemdae Tome", NewJayemdaeTome)
}

// NewJayemdaeTome creates a Jayemdae Tome
// {4} - ARTIFACT
func NewJayemdaeTome(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Jayemdae Tome")
	card.ManaCost = "{4}"
	card.Types = []string{"ARTIFACT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddManaCost("{4}").
		AddTapCost().
		AddEffect(abilities.NewDrawCardsEffect(1)).
		Build()
	card.AddAbility(ability0)
	return card, nil
}
