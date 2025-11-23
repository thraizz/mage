package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/token"
	"github.com/magefree/mage-server-go/internal/game/effects"
)

func init() {
	cards.Register("Galazeth Prismari", NewGalazethPrismari)
}

// NewGalazethPrismari creates a Galazeth Prismari
// {2}{U}{R} - CREATURE
// Flying
func NewGalazethPrismari(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Galazeth Prismari")
	card.ManaCost = "{2}{U}{R}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"ELDER", "DRAGON"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "3"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	token1_0, err := token.GetToken("TreasureToken")
	if err != nil {
		return nil, err
	}
	ability1, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewCreateTokenEffect(token1_0)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability1)
	ability2, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewGrantAbilityEffect(new ConditionalAnyColorManaAbility(1, new InstantOrSorcerySpellManaBuilder()))).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability2)
	return card, nil
}