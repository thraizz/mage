package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Sculptor Of Winter", NewSculptorOfWinter)
}

// NewSculptorOfWinter creates a Sculptor Of Winter
// {1}{G} - CREATURE
func NewSculptorOfWinter(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Sculptor Of Winter")
	card.ManaCost = "{1}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"ELF", "ROGUE"}
	card.Supertypes = []string{"SNOW"}
	card.Power = "2"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddTapCost().
		AddEffect(abilities.NewUntapEffect()).
		Build()
	card.AddAbility(ability0)
	return card, nil
}
