package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Tezzeret Master Of Metal", NewTezzeretMasterOfMetal)
}

// NewTezzeretMasterOfMetal creates a Tezzeret Master Of Metal
// {4}{U}{B} - PLANESWALKER
func NewTezzeretMasterOfMetal(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Tezzeret Master Of Metal")
	card.ManaCost = "{4}{U}{B}"
	card.Types = []string{"PLANESWALKER"}
	card.Subtypes = []string{"TEZZERET"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Loyalty = "5"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewLoseLifeEffect(ArtifactYouControlCount.instance)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
