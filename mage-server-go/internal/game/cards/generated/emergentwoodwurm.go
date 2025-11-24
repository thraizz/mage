package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Emergent Woodwurm", NewEmergentWoodwurm)
}

// NewEmergentWoodwurm creates a Emergent Woodwurm
// {6}{G} - CREATURE
func NewEmergentWoodwurm(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Emergent Woodwurm")
	card.ManaCost = "{6}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"WURM"}
	card.Power = "4"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - LookLibraryAndPickControllerEffect(                 SourcePermanentPowerValue.NOT_NEG...)
	// card.AddAbility(ability0)
	return card, nil
}
