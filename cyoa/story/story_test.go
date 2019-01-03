package story

import "testing"

func TestParseStory(t *testing.T) {
	path := "../data/test_gopher.json"
	story, err := ParseStory(path)

	if err != nil {
		t.Fatal("Couldn't parse JSON file: ", path, "\nError: ", err)
	}

	title := "Test Story"
	paragraphs := []string{"A", "B"}
	options := []Option{{"1", "arc1"}}

	storyName := "intro"
	s := story[storyName]

	if s.Title != title {
		t.Error(
			"Expected title:", title,
			"Received title:", s.Title,
		)
	}

	if len(s.Story) != len(paragraphs) {
		t.Fatal(
			"Number of paragraphs did not match.",
			" Expected: ", len(paragraphs),
			" Received: ", len(s.Story),
		)
	}

	for i, p := range paragraphs {
		if p != s.Story[i] {
			t.Error(
				"Expected paragraph: ", p,
				" Received paragraph: ", s.Story[i],
			)
		}
	}

	if len(options) != len(s.Options) {
		t.Fatal(
			"Number of options did not match.",
			" Expected: ", len(options),
			" Received: ", len(s.Options),
		)
	}

	for i, o := range options {
		if o != s.Options[i] {
			t.Error(
				"Expected option: ", o,
				" Received option: ", s.Options[i],
			)
		}
	}
}
