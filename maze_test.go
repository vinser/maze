package maze

import (
	"testing"
)

func TestNew(t *testing.T) {
	t.Run("Successful creation", func(t *testing.T) {
		m, err := New(41, 21, 11, 7)
		if err != nil {
			t.Fatalf("Expected no error, but got %v", err)
		}
		if m == nil {
			t.Fatal("Expected maze to be created, but got nil")
		}
		if m.Width() != 41 {
			t.Errorf("Expected width 41, got %d", m.Width())
		}
		if m.Height() != 21 {
			t.Errorf("Expected height 21, got %d", m.Height())
		}
		if m.DenWidth() != 11 {
			t.Errorf("Expected den width 11, got %d", m.DenWidth())
		}
		if m.DenHeight() != 7 {
			t.Errorf("Expected den height 7, got %d", m.DenHeight())
		}
	})

	t.Run("Dimension adjustment", func(t *testing.T) {
		m, err := New(40, 20, 10, 6)
		if err != nil {
			t.Fatalf("Expected no error, but got %v", err)
		}
		if m.Width() != 41 {
			t.Errorf("Expected width to be adjusted to 41, got %d", m.Width())
		}
		if m.Height() != 21 {
			t.Errorf("Expected height to be adjusted to 21, got %d", m.Height())
		}
		if m.DenWidth() != 11 {
			t.Errorf("Expected denWidth to be adjusted to 11, got %d", m.denWidth)
		}
		if m.DenHeight() != 7 {
			t.Errorf("Expected denHeight to be adjusted to 7, got %d", m.denHeight)
		}
	})

	t.Run("Invalid dimensions", func(t *testing.T) {
		_, err := New(0, 21, 0, 0)
		if err == nil {
			t.Error("Expected error for zero width, but got nil")
		}
	})

	t.Run("Den too large", func(t *testing.T) {
		_, err := New(21, 21, 21, 5)
		if err == nil {
			t.Error("Expected error for den width too large, but got nil")
		}
	})
}

func TestDenHelpers(t *testing.T) {
	// Create a maze with a known den for consistent testing.
	// Maze is 21x21, den is 7x5.
	// Den starts at ( (21-7)/2, (21-5)/2 ) -> (7, 7) after alignment adjustment.
	// Den is from x=[7, 13], y=[7, 11].
	m, err := New(21, 21, 7, 5)
	if err != nil {
		t.Fatalf("Failed to create maze for testing: %v", err)
	}

	t.Run("IsInsideDen", func(t *testing.T) {
		testCases := []struct {
			name     string
			point    Point
			expected bool
		}{
			{"Center of den", Point{10, 9}, true},
			{"Top-left corner of den", Point{7, 7}, true},
			{"Bottom-right edge of den", Point{13, 11}, true},
			{"Just outside den (top)", Point{10, 6}, false},
			{"Just outside den (left)", Point{6, 9}, false},
			{"Far outside den", Point{1, 1}, false},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				if got := m.IsInsideDen(tc.point); got != tc.expected {
					t.Errorf("IsInsideDen(%+v) = %v; want %v", tc.point, got, tc.expected)
				}
			})
		}
	})

	t.Run("IsAdjacentToDen", func(t *testing.T) {
		testCases := []struct {
			name     string
			point    Point
			expected bool
		}{
			{"Adjacent to top wall", Point{10, 6}, true},
			{"Adjacent to left wall", Point{6, 9}, true},
			{"Adjacent to bottom wall", Point{10, 12}, true},
			{"Adjacent to right wall", Point{14, 9}, true},
			{"Diagonal to den", Point{6, 6}, false},
			{"Inside den", Point{10, 9}, false},
			{"Far from den", Point{1, 1}, false},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				if got := m.IsAdjacentToDen(tc.point); got != tc.expected {
					t.Errorf("IsAdjacentToDen(%+v) = %v; want %v", tc.point, got, tc.expected)
				}
			})
		}
	})
}
