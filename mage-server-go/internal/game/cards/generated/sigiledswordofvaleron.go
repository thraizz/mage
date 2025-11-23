package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/token"
)

func init() {
	cards.Register("Sigiled Sword Of Valeron", NewSigiledSwordOfValeron)
}

// NewSigiledSwordOfValeron creates a Sigiled Sword Of Valeron
// {3} - ARTIFACT
func NewSigiledSwordOfValeron(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Sigiled Sword Of Valeron")
	card.ManaCost = "{3}"
	card.Types = []string{"ARTIFACT"}
	card.Subtypes = []string{"EQUIPMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewGainAbilityAttachedEffect("VigilanceAbility", abilities.AttachmentTypeEquipment)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	token1_0, err := token.GetToken("KnightToken")
	if err != nil {
		return nil, err
	}
	ability1, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewCreateTokenEffect(token1_0, 1, false, true)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability1)
	return card, nil
}
