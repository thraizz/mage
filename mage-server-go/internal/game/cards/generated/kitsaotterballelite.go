package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Kitsa Otterball Elite", NewKitsaOtterballElite)
}

// NewKitsaOtterballElite creates a Kitsa Otterball Elite
// {1}{U} - CREATURE
// Vigilance
func NewKitsaOtterballElite(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Kitsa Otterball Elite")
	card.ManaCost = "{1}{U}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"OTTER", "WIZARD"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "1"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement triggered ability: ActivateIfConditionActivatedAbility
	//   - Effect: CopyTargetStackObjectEffect()
	//
	// Targets:
	//   - abilities.NewTargetRequirement(1, 1, abilities.NewSpellTargetFilter())
	// card.AddAbility(ability0)
	ability1 := abilities.NewKeywordAbility(card.ID, abilities.KeywordVigilance)
	card.AddAbility(ability1)
	return card, nil
}
