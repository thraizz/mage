package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/token"
)

func init() {
	cards.Register("Ashnod Flesh Mechanist", NewAshnodFleshMechanist)
}

// NewAshnodFleshMechanist creates a Ashnod Flesh Mechanist
// {B} - CREATURE
// Deathtouch
func NewAshnodFleshMechanist(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Ashnod Flesh Mechanist")
	card.ManaCost = "{B}"
	card.Types = []string{"CREATURE"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "1"
	card.Toughness = "1"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordDeathtouch)
	card.AddAbility(ability0)
	token1_0, err := token.GetToken("AshnodZombieToken")
	if err != nil {
		return nil, err
	}
	ability1 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddManaCost("{5}").
		AddEffect(abilities.NewCreateTokenEffectAttacking(token1_0, 1, true, false)).
		Build()
	card.AddAbility(ability1)
	// TODO: Implement spell ability with unmapped effects
	//   - DoIfCostPaid(                         new CreateTokenEffect(   ...)
	// card.AddAbility(ability2)
	return card, nil
}
