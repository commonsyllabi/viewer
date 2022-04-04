package commoncartridge

type Cartridge interface {
	// // Load opens an IMSCC archive at the given path
	// Load(string) (Cartridge, error) --> this becomes a package level function
	// AsObject returns a serialized JSON representation - TODO: check how to make the CC implement Marshal
	AsObject() ([]byte, error)
	ParseManifest() (Manifest, error)
	// Dump returns all the contents of the cartridge as a string
	Dump() []string
	// Title returns the title of the loaded cartridge
	Title() string

	// Tobi's weird ideas
	// Weblinks() ([]Weblink, error)
	// GetWeblink(string) (Weblink, error)
	// Webcontents() ([]os.File, error)
	// GetWebcontent(string) (os.File, error)
}
