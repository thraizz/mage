package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/token"
)

func init() {
	cards.Register("Bladewing Deathless Tyrant", NewBladewingDeathlessTyrant)
}

// NewBladewingDeathlessTyrant creates a Bladewing Deathless Tyrant
// {5}{B}{R} - CREATURE
// Flying, Haste
func NewBladewingDeathlessTyrant(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Bladewing Deathless Tyrant")
	card.ManaCost = "{5}{B}{R}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"DRAGON", "SKELETON"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "6"
	card.Toughness = "6"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	ability1 := abilities.NewKeywordAbility(card.ID, abilities.KeywordHaste)
	card.AddAbility(ability1)
	token2_0, err := token.GetToken("ZombieKnightToken")
	if err != nil {
		return nil, err
	}
	ability2, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewCreateTokenEffect(token2_0, new CardsInControllerGraveyardCount())).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability2)
	return card, nil
}