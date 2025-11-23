package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/counters"
)

func init() {
	cards.Register("Astrologians Planisphere", NewAstrologiansPlanisphere)
}

// NewAstrologiansPlanisphere creates a Astrologians Planisphere
// {1}{U} - ARTIFACT
func NewAstrologiansPlanisphere(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Astrologians Planisphere")
	card.ManaCost = "{1}{U}"
	card.Types = []string{"ARTIFACT"}
	card.Subtypes = []string{"EQUIPMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewEquipAbility(card.ID, "{2}", true)
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	ability1, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewAddCountersSourceEffect(counters.CounterTypeP1P1.CreateInstance(1))).
		AddEffect(abilities.NewGainAbilityAttachedEffect(new OrTriggeredAbility( Zone.BATTLEFIELD, new AddCountersSourceEffect(counters.CounterTypeP1P1.CreateInstance(1)), new SpellCastControllerTriggeredAbility(null, false), new DrawNthCardTriggeredAbility(null, false, 3) ), AttachmentType.EQUIPMENT)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability1)
	return card, nil
}