package main

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestReadSentNew(t *testing.T) {
	tmp := t.TempDir()
	filename := filepath.Join(tmp, "does-not-exist.json")
	got, err := readSent(filename)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) > 0 {
		t.Errorf("Expect empty slice, but got: %v", got)
	}
}

func TestReadSent(t *testing.T) {
	tmp := t.TempDir()
	filename := filepath.Join(tmp, "sent.json")
	data := []byte(`[{"id":156522206,"price":"£2,500"}]`)
	os.WriteFile(filename, data, 0644)

	got, err := readSent(filename)
	if err != nil {
		t.Fatal(err)
	}
	want := []flat{
		flat{ID: 156522206, Price: "£2,500"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("readSent failed: got: %v, want: %v", got, want)
	}
}

func TestRemoveAlreadySent(t *testing.T) {
	flats := []flat{
		flat{
			ID:    156522206,
			Price: "£2,500",
		},
		flat{
			ID:    158462822,
			Price: "£3,000",
		}}
	sent := []flat{
		flat{
			ID:    156522206,
			Price: "£2,500",
		},
	}

	got := removeAlreadySent(flats, sent)
	want := []flat{flat{ID: 158462822, Price: "£3,000"}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("removeAlreadySent failed: got: %v, want: %v",
			got, want)
	}

}

func TestRemoveDelisted(t *testing.T) {
	flats := []flat{
		flat{
			ID:    156522206,
			Price: "£2,500",
		},
		flat{
			ID:    158462822,
			Price: "£3,000",
		}}
	sent := []flat{
		flat{
			ID:    156522206,
			Price: "£2,500",
		},
		flat{
			ID:    157948184,
			Price: "£2,400",
		},
	}

	got := removeDelisted(sent, flats)
	want := []flat{flat{ID: 156522206, Price: "£2,500"}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("removeDelisted failed: got: %v, want: %v",
			got, want)
	}
}

func TestWriteSentNew(t *testing.T) {
	tmp := t.TempDir()
	filename := filepath.Join(tmp, "sent.json")
	flats := []flat{flat{ID: 156522206, Price: "£2,500"}}

	err := writeSent(flats, filename)
	if err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(filename)
	if err != nil {
		t.Fatal(err)
	}
	want := []byte(`[{"id":156522206,"price":"£2,500"}]`)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("writeSent failed: got: %v, want: %v", got, want)
	}
}

func TestWriteSentOverride(t *testing.T) {
	tmp := t.TempDir()
	filename := filepath.Join(tmp, "sent.json")
	_, err := os.Create(filename)
	if err != nil {
		t.Fatal(err)
	}
	flats := []flat{flat{ID: 156522206, Price: "£2,500"}}
	err = writeSent(flats, filename)

	if err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(filename)
	if err != nil {
		t.Fatal(err)
	}
	want := []byte(`[{"id":156522206,"price":"£2,500"}]`)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("writeSent failed: got: %v, want: %v", got, want)
	}
}
