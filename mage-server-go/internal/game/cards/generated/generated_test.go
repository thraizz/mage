package generated_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game/cards"
	_ "github.com/magefree/mage-server-go/internal/game/cards/generated" // Import to register cards
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGeneratedCards(t *testing.T) {
	ownerID := uuid.New()

	testCases := []struct {
		name              string
		expectedCost      string
		expectedTypes     []string
		expectedP         string
		expectedT         string
		expectedAbilities int
	}{
		{
			name:              "Lightning Bolt",
			expectedCost:      "{R}",
			expectedTypes:     []string{"INSTANT"},
			expectedAbilities: 1,
		},
		{
			name:              "Llanowar Elves",
			expectedCost:      "{G}",
			expectedTypes:     []string{"CREATURE"},
			expectedP:         "1",
			expectedT:         "1",
			expectedAbilities: 1,
		},
		{
			name:              "Serra Angel",
			expectedCost:      "{3}{W}{W}",
			expectedTypes:     []string{"CREATURE"},
			expectedP:         "4",
			expectedT:         "4",
			expectedAbilities: 2, // Flying + Vigilance
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			builder, ok := cards.Registry.Get(tc.name)
			require.True(t, ok, "%s should be registered", tc.name)

			card, err := builder(ownerID, nil)
			require.NoError(t, err)
			assert.Equal(t, tc.name, card.Name)
			assert.Equal(t, tc.expectedCost, card.ManaCost)
			assert.Equal(t, tc.expectedTypes, card.Types)

			if tc.expectedP != "" {
				assert.Equal(t, tc.expectedP, card.Power)
			}
			if tc.expectedT != "" {
				assert.Equal(t, tc.expectedT, card.Toughness)
			}

			assert.Len(t, card.Abilities, tc.expectedAbilities,
				"%s should have %d abilities", tc.name, tc.expectedAbilities)
		})
	}
}
