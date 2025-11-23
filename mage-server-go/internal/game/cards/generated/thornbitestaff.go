package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Thornbite Staff", NewThornbiteStaff)
}

// NewThornbiteStaff creates a Thornbite Staff
// {2} - KINDRED ARTIFACT
func NewThornbiteStaff(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Thornbite Staff")
	card.ManaCost = "{2}"
	card.Types = []string{"KINDRED", "ARTIFACT"}
	card.Subtypes = []string{"SHAMAN", "EQUIPMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewGainAbilityAttachedEffect(gainedAbility, AttachmentType.EQUIPMENT)).
		AddEffect(abilities.NewUntapEffect()).
		AddEffect(abilities.NewGainAbilityAttachedEffect(new DiesCreatureTriggeredAbility(new UntapSourceEffect(),false), AttachmentType.EQUIPMENT)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	ability1 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddManaCost("{2}").
		AddTapCost().
		AddEffect(abilities.NewDamageEffect(1)).
		Build()
	card.AddAbility(ability1)
	ability2, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewAttachEffect(abilities.OutcomeDetriment)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability2)
	return card, nil
}