package main

const (
	spanish   = "Spanish"
	french    = "French"
	norwegian = "Norwegian"

	englishHelloPrefix   = "Hello, "
	spanishHelloPrefix   = "Hola, "
	frenchHelloPrefix    = "Bonjour, "
	norwegianHelloPrefix = "Hei, "
)

func Hello(name string, lang string) string {
	if name == "" {
		name = "World"
	}

	return greetingPrefix(lang) + name
}

func greetingPrefix(lang string) (prefix string) {
	switch lang {
	case spanish:
		prefix = spanishHelloPrefix
	case french:
		prefix = frenchHelloPrefix
	case norwegian:
		prefix = norwegianHelloPrefix
	default:
		prefix = englishHelloPrefix
	}
	return
}
