package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/effects"
	"github.com/magefree/mage-server-go/internal/game/token"
)

func init() {
	cards.Register("Giants Amulet", NewGiantsAmulet)
}

// NewGiantsAmulet creates a Giants Amulet
// {U} - ARTIFACT
func NewGiantsAmulet(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Giants Amulet")
	card.ManaCost = "{U}"
	card.Types = []string{"ARTIFACT"}
	card.Subtypes = []string{"EQUIPMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewBoostEquippedEffect(0, 1)).
		AddEffect(abilities.NewGrantAbilityEffect("HexproofAbility", effects.DurationPermanent)).
		AddEffect(abilities.NewGainAbilityAttachedEffect(abilities.NewKeywordAbility(card.ID, abilities.KeywordHexproof), abilities.AttachmentTypeEquipment, abilities.DurationWhileOnBattlefield, "")).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	// TODO: Implement spell ability with unmapped effects
	//   - DoIfCostPaid(                 new CreateTokenAttachSourceEffect...)
	// card.AddAbility(ability1)
	return card, nil
}
