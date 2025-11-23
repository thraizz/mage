package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Thornvault Forager", NewThornvaultForager)
}

// NewThornvaultForager creates a Thornvault Forager
// {1}{G} - CREATURE
func NewThornvaultForager(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Thornvault Forager")
	card.ManaCost = "{1}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"SQUIRREL", "RANGER"}
	card.Power = "2"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.BuildSimpleManaAbility(card.ID, "G")
	card.AddAbility(ability0)
	ability1 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddTapCost().
		AddEffect(abilities.NewSearchLibraryPutInHandEffect(abilities.NewTargetRequirement(0, 1, abilities.NewAnyTargetFilter()), true)).
		Build()
	card.AddAbility(ability1)
	return card, nil
}
