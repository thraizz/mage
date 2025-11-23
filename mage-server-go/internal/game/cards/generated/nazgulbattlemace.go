package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Nazgul Battle Mace", NewNazgulBattleMace)
}

// NewNazgulBattleMace creates a Nazgul Battle Mace
// {5} - ARTIFACT
func NewNazgulBattleMace(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Nazgul Battle Mace")
	card.ManaCost = "{5}"
	card.Types = []string{"ARTIFACT"}
	card.Subtypes = []string{"EQUIPMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewEquipAbility(card.ID, "{3}", false)
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	ability1, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewGainAbilityAttachedEffect(new MenaceAbility(false), AttachmentType.EQUIPMENT)).
		AddEffect(abilities.NewGainAbilityAttachedEffect("DeathtouchAbility", abilities.AttachmentTypeEquipment)).
		AddEffect(abilities.NewGainAbilityAttachedEffect(new AnnihilatorAbility(1), AttachmentType.EQUIPMENT)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability1)
	return card, nil
}