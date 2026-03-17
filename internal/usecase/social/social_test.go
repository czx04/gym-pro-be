package social

import (
	"reflect"
	"testing"
)

func TestNormalizeHashtags(t *testing.T) {
	caption := "Hom nay leg day #LegDay #fitness #fitness!"

	tags, err := normalizeHashtags(&caption, []string{"  #GymLife ", "fitness", ""})
	if err != nil {
		t.Fatalf("normalizeHashtags returned error: %v", err)
	}

	want := []string{"gymlife", "fitness", "legday"}
	if !reflect.DeepEqual(tags, want) {
		t.Fatalf("unexpected hashtags: got %v want %v", tags, want)
	}
}

func TestNormalizeContentTypeDefaultsToGeneral(t *testing.T) {
	if got := normalizeContentType(nil); got != "general" {
		t.Fatalf("expected default content type to be general, got %q", got)
	}
}
