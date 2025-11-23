package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Morbius The Living Vampire", NewMorbiusTheLivingVampire)
}

// NewMorbiusTheLivingVampire creates a Morbius The Living Vampire
// {2}{U}{B} - CREATURE
// Flying, Vigilance, Lifelink
func NewMorbiusTheLivingVampire(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Morbius The Living Vampire")
	card.ManaCost = "{2}{U}{B}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"VAMPIRE", "SCIENTIST", "VILLAIN"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "3"
	card.Toughness = "1"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	ability1 := abilities.NewKeywordAbility(card.ID, abilities.KeywordVigilance)
	card.AddAbility(ability1)
	ability2 := abilities.NewKeywordAbility(card.ID, abilities.KeywordLifelink)
	card.AddAbility(ability2)
	// TODO: Implement activated ability with unmapped effects
	//   - LookLibraryAndPickControllerEffect(3, 1, PutCards.HAND, PutCards.BOTTOM_ANY)
	// card.AddAbility(ability3)
	return card, nil
}
