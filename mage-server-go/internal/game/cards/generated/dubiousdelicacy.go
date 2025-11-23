package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Dubious Delicacy", NewDubiousDelicacy)
}

// NewDubiousDelicacy creates a Dubious Delicacy
// {2}{B} - ARTIFACT
// Flash
func NewDubiousDelicacy(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Dubious Delicacy")
	card.ManaCost = "{2}{B}"
	card.Types = []string{"ARTIFACT"}
	card.Subtypes = []string{"FOOD"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlash)
	card.AddAbility(ability0)
	ability1, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewBoostEffect(-3, -3)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability1)
	ability2 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddManaCost("{2}").
		AddTapCost().
		AddSacrificeSourceCost().
		AddEffect(abilities.NewLoseLifeEffect(3)).
		Build()
	card.AddAbility(ability2)
	return card, nil
}
