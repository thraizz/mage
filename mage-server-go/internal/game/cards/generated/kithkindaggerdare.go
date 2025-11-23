package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Kithkin Daggerdare", NewKithkinDaggerdare)
}

// NewKithkinDaggerdare creates a Kithkin Daggerdare
// {1}{G} - CREATURE
func NewKithkinDaggerdare(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Kithkin Daggerdare")
	card.ManaCost = "{1}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"KITHKIN", "SOLDIER"}
	card.Power = "1"
	card.Toughness = "1"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddTapCost().
		AddEffect(abilities.NewBoostEffect(2, 2)).
		Build()
	card.AddAbility(ability0)
	return card, nil
}
