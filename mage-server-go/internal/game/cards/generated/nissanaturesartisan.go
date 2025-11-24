package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Nissa Natures Artisan", NewNissaNaturesArtisan)
}

// NewNissaNaturesArtisan creates a Nissa Natures Artisan
// {4}{G}{G} - PLANESWALKER
func NewNissaNaturesArtisan(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Nissa Natures Artisan")
	card.ManaCost = "{4}{G}{G}"
	card.Types = []string{"PLANESWALKER"}
	card.Subtypes = []string{"NISSA"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Loyalty = "5"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewGainLifeEffect(3)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
