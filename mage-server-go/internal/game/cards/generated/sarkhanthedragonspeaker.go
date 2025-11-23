package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Sarkhan The Dragonspeaker", NewSarkhanTheDragonspeaker)
}

// NewSarkhanTheDragonspeaker creates a Sarkhan The Dragonspeaker
// {3}{R}{R} - PLANESWALKER
// Flying, Indestructible, Haste
func NewSarkhanTheDragonspeaker(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Sarkhan The Dragonspeaker")
	card.ManaCost = "{3}{R}{R}"
	card.Types = []string{"PLANESWALKER"}
	card.Subtypes = []string{"SARKHAN"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Loyalty = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	ability1 := abilities.NewKeywordAbility(card.ID, abilities.KeywordIndestructible)
	card.AddAbility(ability1)
	ability2 := abilities.NewKeywordAbility(card.ID, abilities.KeywordHaste)
	card.AddAbility(ability2)
	ability3, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewDamageEffect(4)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability3)
	return card, nil
}
