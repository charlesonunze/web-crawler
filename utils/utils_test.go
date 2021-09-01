package utils

import (
	"reflect"
	"testing"

	"golang.org/x/net/html"
)

func TestGetHTMLPage(t *testing.T) {
	t.Run("GetHTMLPage", func(t *testing.T) {
		page, err := GetHTMLPage("https://motherfuckingwebsite.com")
		p := &html.Node{}
		if err != nil {
			t.Error("Error ->", err)
		}

		if reflect.TypeOf(page.Type) != reflect.TypeOf(p.Type) {
			t.Errorf("not an html page!")
		}
	})
}

func TestExtractLinks(t *testing.T) {
	t.Run("ExtractLinks", func(t *testing.T) {
		page, err := GetHTMLPage("https://motherfuckingwebsite.com")
		if err != nil {
			t.Error("Error ->", err)
		}

		links := ExtractLinks(nil, page)

		page, err = GetHTMLPage(links[0])
		if err != nil {
			t.Error("error ->", err)
		}
	})
}

func TestToWWW(t *testing.T) {
	got := ToWWW("https://motherfuckingwebsite.com")
	want := "https://www.motherfuckingwebsite.com"

	if got != want {
		t.Errorf("want %s but got %s", want, got)
	}
}

func TestBelongsToSubdomain(t *testing.T) {
	got := BelongsToSubdomain("/careers", "https://monzo.com")
	want := true

	if got != want {
		t.Errorf("want %v but got %v", want, got)
	}
}

func TestRemoveTrailingSlash(t *testing.T) {
	got := RemoveTrailingSlash("/")
	want := ""

	if got != want {
		t.Errorf("want %s but got %s", want, got)
	}
}

func TestFormatURL(t *testing.T) {
	got := FormatURL("/blog", "https://monzo")
	want := "https://monzo/blog"

	if got != want {
		t.Errorf("want %s but got %s", want, got)
	}
}
