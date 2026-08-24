package models

// Reading holds how a dhikr/salavat preset is read, shown via the goal
// card's "Okunuş" button. ImagePath, when set, points at a static scan of
// the preset's page (Arabic text, Turkish transliteration and translation
// all together) and takes priority over the plain Arabic/Turkish fields.
// Steps, when set, breaks the reading into an ordered sequence of distinct
// phrases (each read its own number of times) instead of one fixed text —
// e.g. Cünnetü'l-Esmâ, which is read differently for its first 11
// repetitions than its last 3.
type Reading struct {
	Arabic    string
	Turkish   string
	ImagePath string
	Steps     []ReadingStep
	Note      string
}

// ReadingStep is one phrase within a multi-step Reading, read Repeat times
// before moving on to the next step.
type ReadingStep struct {
	Arabic  string
	Turkish string
	Repeat  int
}

// IsEmpty reports whether nothing has been filled in yet for this reading.
func (r Reading) IsEmpty() bool {
	return r.Arabic == "" && r.Turkish == "" && r.ImagePath == "" && len(r.Steps) == 0
}

// TitlePresets are the choices offered in the "Münacaat" dropdown.
var TitlePresets = []string{
	"Salât-ı Münciye",
	"Salât-ı Fethiyye",
	"Salât-ı Nâriye",
	"Salavât-ı Şerife",
	"Kelime-i Tevhid",
	"İstiğfar-ı Şerif",
	"Havkale (Lâ Havle)",
	"Cünnetü'l-Esmâ",
}

// titleReadings maps a preset to its Reading. Left blank for now — to be
// filled in later.
var titleReadings = map[string]Reading{
	"Salât-ı Münciye":  {ImagePath: "/static/images/readings/munciye.png"},
	"Salât-ı Fethiyye": {ImagePath: "/static/images/readings/fethiyye.png"},
	"Salât-ı Nâriye":   {ImagePath: "/static/images/readings/nariye.png"},
	"İstiğfar-ı Şerif": {ImagePath: "/static/images/readings/istigfar.jpeg"},
	"Havkale (Lâ Havle)": {
		Arabic:  "لَا حَوْلَ وَلَا قُوَّةَ إِلَّا بِاللَّهِ الْعَلِيِّ الْعَظِيمِ",
		Turkish: "Lâ havle velâ kuvvete illâ billâhil aliyyil azîm.",
	},
	"Salavât-ı Şerife": {
		Arabic:  "اللَّهُمَّ صَلِّ عَلَى سَيِّدِنَا مُحَمَّدٍ وَعَلَى آلِ سَيِّدِنَا مُحَمَّدٍ",
		Turkish: "Allâhümme salli alâ seyyidinâ Muhammedin ve alâ âli seyyidinâ Muhammed.",
	},
	"Kelime-i Tevhid": {
		Arabic:  "لَا إِلَٰهَ إِلَّا اللَّهُ مُحَمَّدٌ رَسُولُ اللَّهِ",
		Turkish: "Lâ ilâhe illallah, Muhammedün Rasûlullah.",
	},
	"Cünnetü'l-Esmâ": {
		Steps: []ReadingStep{
			{
				Arabic:  "بِسْمِ اللَّهِ الرَّحْمَٰنِ الرَّحِيمِ\nفَرْدٌ حَيٌّ قَيُّومٌ حَكَمٌ عَدْلٌ قُدُّوسٌ\nعَنَتِ الْوُجُوهُ لِلْحَيِّ الْقَيُّومِ",
				Turkish: "Bismillâhirrahmânirrahîm.\nFerdün, Hayyun, Kayyûmun, Hakemün, Adlün, Kuddûsün.\nAnetil vücûhu lil-hayyil-kayyûm.",
				Repeat:  11,
			},
			{
				Arabic:  "بِسْمِ اللَّهِ الرَّحْمَٰنِ الرَّحِيمِ\nفَرْدٌ حَيٌّ قَيُّومٌ حَكَمٌ عَدْلٌ قُدُّوسٌ\nنَجِّنِي مِنَ الْقَوْمِ الظَّالِمِينَ",
				Turkish: "Bismillâhirrahmânirrahîm.\nFerdün, Hayyun, Kayyûmun, Hakemün, Adlün, Kuddûsün.\nNeccinî minel-kavmiz-zâlimîn.",
				Repeat:  3,
			},
		},
	},
}

// ReadingFor looks up the Reading for a goal title (a preset from
// TitlePresets). ok is false if title isn't a known preset.
func ReadingFor(title string) (reading Reading, ok bool) {
	reading, ok = titleReadings[title]
	return reading, ok
}
