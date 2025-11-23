package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/token"
)

func init() {
	cards.Register("Alela Artful Provocateur", NewAlelaArtfulProvocateur)
}

// NewAlelaArtfulProvocateur creates a Alela Artful Provocateur
// {1}{W}{U}{B} - CREATURE
// Flying, Deathtouch, Lifelink
func NewAlelaArtfulProvocateur(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Alela Artful Provocateur")
	card.ManaCost = "{1}{W}{U}{B}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"FAERIE", "WARLOCK"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "2"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	ability1 := abilities.NewKeywordAbility(card.ID, abilities.KeywordDeathtouch)
	card.AddAbility(ability1)
	ability2 := abilities.NewKeywordAbility(card.ID, abilities.KeywordLifelink)
	card.AddAbility(ability2)
	ability3, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewBoostEffect(1, 0, filter, true)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability3)
	token4_0, err := token.GetToken("FaerieToken")
	if err != nil {
		return nil, err
	}
	ability4, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewCreateTokenEffect(token4_0)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability4)
	return card, nil
}
