package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Avarice Amulet", NewAvariceAmulet)
}

// NewAvariceAmulet creates a Avarice Amulet
// {4} - ARTIFACT
func NewAvariceAmulet(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Avarice Amulet")
	card.ManaCost = "{4}"
	card.Types = []string{"ARTIFACT"}
	card.Subtypes = []string{"EQUIPMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewBoostEquippedEffect(2, 0)).
		AddEffect(abilities.NewGainAbilityAttachedEffect("VigilanceAbility", abilities.AttachmentTypeEquipment)).
		AddEffect(abilities.NewDrawCardsEffect(1)).
		AddEffect(abilities.NewGainAbilityAttachedEffect(new BeginningOfUpkeepTriggeredAbility(new DrawCardSourceControllerEffect(1)), AttachmentType.EQUIPMENT)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}