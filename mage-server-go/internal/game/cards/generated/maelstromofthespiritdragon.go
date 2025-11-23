package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Maelstrom Of The Spirit Dragon", NewMaelstromOfTheSpiritDragon)
}

// NewMaelstromOfTheSpiritDragon creates a Maelstrom Of The Spirit Dragon
//   - LAND
func NewMaelstromOfTheSpiritDragon(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Maelstrom Of The Spirit Dragon")
	card.ManaCost = ""
	card.Types = []string{"LAND"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.BuildSimpleManaAbility(card.ID, "C")
	card.AddAbility(ability0)
	ability1 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddManaCost("{4}").
		AddTapCost().
		AddSacrificeSourceCost().
		AddEffect(abilities.NewSearchLibraryPutInHandEffect(abilities.NewTargetRequirement(0, 1, abilities.NewAnyTargetFilter()), true)).
		Build()
	card.AddAbility(ability1)
	return card, nil
}
