package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Casey Jones Jury Rig Justiciar", NewCaseyJonesJuryRigJusticiar)
}

// NewCaseyJonesJuryRigJusticiar creates a Casey Jones Jury Rig Justiciar
// {1}{R} - CREATURE
// Haste
func NewCaseyJonesJuryRigJusticiar(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Casey Jones Jury Rig Justiciar")
	card.ManaCost = "{1}{R}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"HUMAN", "BERSERKER"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "2"
	card.Toughness = "1"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordHaste)
	card.AddAbility(ability0)
	// TODO: Implement spell ability with unmapped effects
	//   - LookLibraryAndPickControllerEffect(                 4, 1, StaticFilters.FILTER_CARD_A...)
	// card.AddAbility(ability1)
	return card, nil
}
