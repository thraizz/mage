package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Ezrim Agency Chief", NewEzrimAgencyChief)
}

// NewEzrimAgencyChief creates a Ezrim Agency Chief
// {1}{W}{W}{U}{U} - CREATURE
// Flying, Vigilance, Lifelink, Hexproof
func NewEzrimAgencyChief(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Ezrim Agency Chief")
	card.ManaCost = "{1}{W}{W}{U}{U}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"ARCHON", "DETECTIVE"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "5"
	card.Toughness = "5"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	ability1 := abilities.NewKeywordAbility(card.ID, abilities.KeywordVigilance)
	card.AddAbility(ability1)
	ability2 := abilities.NewKeywordAbility(card.ID, abilities.KeywordLifelink)
	card.AddAbility(ability2)
	ability3 := abilities.NewKeywordAbility(card.ID, abilities.KeywordHexproof)
	card.AddAbility(ability3)
	// TODO: Implement activated ability with unmapped effects
	//   - GainsChoiceOfAbilitiesEffect()
	//
	// Costs:
	//   - AddManaCost("{1}")
	// card.AddAbility(ability4)
	return card, nil
}
