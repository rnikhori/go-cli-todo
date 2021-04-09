package main

import (
	"reflect"
	"testing"
)

func TestNormalizeTags(t *testing.T) {
	cases := []struct {
		name string
		in   []string
		want []string
	}{
		{"empty", nil, nil},
		{"trims and lowercases", []string{" Work ", "URGENT"}, []string{"work", "urgent"}},
		{"dedupes", []string{"work", "Work", "work"}, []string{"work"}},
		{"drops empty", []string{"", "  ", "home"}, []string{"home"}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := normalizeTags(c.in)
			if !reflect.DeepEqual(got, c.want) {
				t.Errorf("normalizeTags(%v) = %v, want %v", c.in, got, c.want)
			}
		})
	}
}

func TestHasTag(t *testing.T) {
	tags := []string{"work", "urgent"}
	if !hasTag(tags, "work") {
		t.Error("expected hasTag to find exact match")
	}
	if !hasTag(tags, " WORK ") {
		t.Error("expected hasTag to be case-insensitive and trim whitespace")
	}
	if hasTag(tags, "home") {
		t.Error("expected hasTag to return false for missing tag")
	}
	if hasTag(nil, "work") {
		t.Error("expected hasTag to return false for nil tags")
	}
}
