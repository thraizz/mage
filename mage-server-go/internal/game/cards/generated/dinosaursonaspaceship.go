package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/effects"
	"github.com/magefree/mage-server-go/internal/game/token"
)

func init() {
	cards.Register("Dinosaurs On A Spaceship", NewDinosaursOnASpaceship)
}

// NewDinosaursOnASpaceship creates a Dinosaurs On A Spaceship
// {4}{R}{W} - CREATURE
// Vigilance, Trample
func NewDinosaursOnASpaceship(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Dinosaurs On A Spaceship")
	card.ManaCost = "{4}{R}{W}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"DINOSAUR"}
	card.Power = "7"
	card.Toughness = "7"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordVigilance)
	card.AddAbility(ability0)
	ability1 := abilities.NewKeywordAbility(card.ID, abilities.KeywordTrample)
	card.AddAbility(ability1)
	ability2, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewBoostEffect(1, 1, filterCreature, true)).
		AddEffect(abilities.NewGrantAbilityEffect("VigilanceAbility", effects.DurationWhileOnBattlefield)).
		AddEffect(abilities.NewGrantAbilityEffect("TrampleAbility", effects.DurationWhileOnBattlefield)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability2)
	token3_0, err := token.GetToken("DinosaurFlyingHasteToken")
	if err != nil {
		return nil, err
	}
	ability3, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewCreateTokenEffect(token3_0)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability3)
	return card, nil
}
