package maps

import "testing"

func TestSearch(t *testing.T) {
	dictionary := Dictionary{"test": "this is just a test"}

	t.Run("known word", func(t *testing.T) {
		word := "test"
		got, _ := dictionary.Search(word)
		want := "this is just a test"
		assertStrings(t, got, want, word)
	})

	t.Run("unknown word", func(t *testing.T) {
		word := "unknown"
		_, err := dictionary.Search(word)

		if err == nil {
			t.Fatal("expected to get an error")
		}

		assertError(t, err, ErrNotFound, word)
	})
}

func TestAdd(t *testing.T) {

	t.Run("add word", func(t *testing.T) {
		dictionary := Dictionary{}
		word := "test"
		def := "this is just a test"
		err := dictionary.Add(word, def)

		assertError(t, err, nil, word)
		assertDef(t, dictionary, word, def)
	})

	t.Run("add existing word", func(t *testing.T) {
		word := "existing"
		def := "this is a existing def"
		dictionary := Dictionary{word: def}
		err := dictionary.Add(word, def)

		if err == nil {
			t.Fatal("expected to get an error")
		}

		assertError(t, err, ErrWordExists, word)
		assertDef(t, dictionary, word, def)
	})
}

func TestUpdate(t *testing.T) {
	word := "test"
	def := "this is just a test"
	newDef := "new definition"
	t.Run("existing word", func(t *testing.T) {
		dictionary := Dictionary{word: def}
		err := dictionary.Update(word, newDef)

		assertError(t, err, nil, word)
		assertDef(t, dictionary, word, newDef)
	})

	t.Run("new word", func(t *testing.T) {
		dictionary := Dictionary{}
		err := dictionary.Update(word, newDef)

		assertError(t, err, ErrNotFound, word)
		assertDef(t, dictionary, word, newDef)
	})
}

func TestDelete(t *testing.T) {
	word := "test"
	def := "this is a test"

	t.Run("word exist", func(t *testing.T) {
		dictionary := Dictionary{word: def}
		err := dictionary.Delete(word)

		assertError(t, err, nil, word)

		_, err = dictionary.Search(word)

		assertError(t, err, ErrNotFound, word)
	})

	t.Run("word doesn't exist", func(t *testing.T) {
		dictionary := Dictionary{}
		err := dictionary.Delete(word)

		assertError(t, err, ErrWordDoesNotExist, word)
	})
}

func assertStrings(t testing.TB, got, want, word string) {
	t.Helper()

	if got != want {
		t.Errorf("got %q want %q given, %s", got, want, word)
	}
}

func assertError(t testing.TB, got, want error, word string) {
	t.Helper()

	if got != want {
		t.Errorf("got %q want %q given, %s", got, want, word)

	}
}

func assertDef(t testing.TB, dictionary Dictionary, word, def string) {
	t.Helper()

	got, err := dictionary.Search(word)
	if err != nil {
		t.Fatal("should find added word:", err)
	}

	assertStrings(t, got, def, word)

}
