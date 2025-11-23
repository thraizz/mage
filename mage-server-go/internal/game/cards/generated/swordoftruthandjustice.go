package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Sword Of Truth And Justice", NewSwordOfTruthAndJustice)
}

// NewSwordOfTruthAndJustice creates a Sword Of Truth And Justice
// {3} - ARTIFACT
func NewSwordOfTruthAndJustice(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Sword Of Truth And Justice")
	card.ManaCost = "{3}"
	card.Types = []string{"ARTIFACT"}
	card.Subtypes = []string{"EQUIPMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewBoostEquippedEffect(2, 2)).
		AddEffect(abilities.NewGainAbilityAttachedEffect(ProtectionAbility.from(ObjectColor.WHITE, ObjectColor.BLUE), AttachmentType.EQUIPMENT)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
