package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Alquist Proft Master Sleuth", NewAlquistProftMasterSleuth)
}

// NewAlquistProftMasterSleuth creates a Alquist Proft Master Sleuth
// {1}{W}{U} - CREATURE
// Vigilance
func NewAlquistProftMasterSleuth(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Alquist Proft Master Sleuth")
	card.ManaCost = "{1}{W}{U}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"HUMAN", "DETECTIVE"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "3"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordVigilance)
	card.AddAbility(ability0)
	ability1 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddTapCost().
		AddEffect(abilities.NewDrawCardsEffect(GetXValue.instance, true)).
		AddEffect(abilities.NewGainLifeEffect(GetXValue.instance)).
		Build()
	card.AddAbility(ability1)
	return card, nil
}