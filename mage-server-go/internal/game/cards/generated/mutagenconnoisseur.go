package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Mutagen Connoisseur", NewMutagenConnoisseur)
}

// NewMutagenConnoisseur creates a Mutagen Connoisseur
// {1}{G}{U} - CREATURE
// Flying, Vigilance
func NewMutagenConnoisseur(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Mutagen Connoisseur")
	card.ManaCost = "{1}{G}{U}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"VEDALKEN", "MUTANT"}
	card.Power = "0"
	card.Toughness = "5"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	ability1 := abilities.NewKeywordAbility(card.ID, abilities.KeywordVigilance)
	card.AddAbility(ability1)
	ability2, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewBoostEffect(xValue, StaticValue.get(0))).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability2)
	return card, nil
}
