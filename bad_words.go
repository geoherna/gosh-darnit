package goshdarnit

// wordGroup is a set of listed words that share a severity.
type wordGroup struct {
	severity Severity
	words    []string
}

// wordGroups lists profanities, slurs, and vulgarities for content moderation,
// grouped by category and then by severity. A word's severity rates how
// offensive it is in its usual sense: "queen", "trap" and "cracker" are mild
// because they are ordinary words far more often than slurs. Every listed
// word goes into the Aho-Corasick automaton at package initialization time.
var wordGroups = []wordGroup{
	// General profanity
	{SeverityModerate, []string{
		"fuck", "fucking", "fucked", "fucker", "fuckers", "fucks", "fuckhead", "fuckface",
		"fuckwit", "fuckwad", "fuckboy", "fuckboi", "fuckup", "fuckups",
		"motherfuck", "motherfucker", "motherfucking", "motherfuckers",
		"shit", "shits", "shitty", "shitting", "shitted", "shitstain", "shithead",
		"shitheads", "shitface", "shitfaced", "shithole", "shithouse", "shitbag",
		"bullshit", "bullshitting", "horseshit", "dipshit", "dumbshit", "apeshit",
		"ass", "asses", "asshole", "assholes", "asshat", "asswipe", "assface",
		"assclown", "jackass", "dumbass", "fatass", "lardass", "smartass", "hardass",
		"badass", "kickass", "candyass", "halfass", "halfassed",
		"bitch", "bitches", "bitchy", "bitching", "bitched", "sonofabitch", "sonuvabitch",
		"bastard", "bastards",
	}},
	{SeverityMild, []string{
		"damn", "damned", "damnit", "goddamn", "goddamned", "goddamnit",
		"hell", "hellhole",
		"crap", "crappy", "crapping",
		"piss", "pissed", "pisses", "pissing", "pissoff",
	}},
	{SeveritySevere, []string{
		"fucktard",
		"whore", "whores", "whorish", "whorey",
		"slut", "sluts", "slutty", "slutbag",
		"skank", "skanky", "skanks",
	}},

	// Sexual terms
	{SeveritySevere, []string{
		"cocksucker", "cocksucking", "cocksuckers",
		"cunt", "cunts", "cunty",
		"titfuck",
		"cumshot", "cumslut", "cumdumpster",
		"blowjob", "blowjobs", "bj",
		"handjob", "handjobs",
		"dildo", "dildos",
		"buttfuck", "buttfucker", "buttfucking",
		"assfuck", "assfucker", "assfucking",
		"porn", "porno",
	}},
	{SeverityModerate, []string{
		"cock", "cocks", "cockhead", "cockblock", "cockface", "cockwomble",
		"dick", "dicks", "dickhead", "dickheads", "dickface", "dickwad", "dickweed",
		"dicknose", "dickbag", "dickless",
		"pussy", "pussies",
		"twat", "twats", "twatwad",
		"tit", "tits", "titty", "titties",
		"cum", "cums", "cumming", "cumstain",
		"jizz", "jizzed", "jizzing",
		"wank", "wanker", "wankers", "wanking",
		"jerkoff", "jackoff",
		"anal",
		"fellatio", "cunnilingus",
		"masturbate", "masturbation", "masturbating",
		"orgasm", "orgasmic",
		"boner", "boners", "hardon", "hardons",
		"horny", "horndog",
		// Lease use clauses name these ("no pornographic materials").
		"pornography", "pornographic",
		"smut", "smutty",
		"hooker", "hookers",
		"pimp", "pimps", "pimping",
	}},
	{SeverityMild, []string{
		"penis", "penises",
		"vagina", "vaginas", "vag", "vajayjay",
		"pussycat",
		"boob", "boobs", "booby", "boobies",
		"erection",
		"hump", "humping", "humped",
		"screw", "screwing", "screwed",
		"bang", "banging", "banged",
		"prostitute", "prostitutes",
		"brothel", "brothels",
	}},

	// Racial slurs
	{SeveritySevere, []string{
		"nigger", "niggers", "nigga", "niggas",
		"chink", "chinks", "chinky",
		"gook", "gooks",
		"spic", "spics", "spick", "spicks",
		"wetback", "wetbacks",
		"beaner", "beaners",
		"kike", "kikes",
		"yid", "yids",
		"heeb", "heebs",
		"towelhead", "towelheads",
		"sandnigger", "sandniggers",
		"camel jockey", "cameljockey",
		"raghead", "ragheads",
		"curry muncher",
		"wog", "wogs",
		"abo", "abos",
		"coon", "coons",
		"darkie", "darkies", "darky",
		"redskin", "redskins",
		"injun", "injuns",
		"squaw", "squaws",
		"zipperhead", "zipperheads",
		"slopehead",
		"jigaboo", "jigaboos", "jiggaboo",
		"sambo", "sambos",
		"spearchucker", "spearchuckers",
		"porchmonkey", "porch monkey",
		"pickaninny", "picaninny",
	}},
	{SeverityModerate, []string{
		// Also ordinary words: Spanish "negro", the Negro Leagues, "Jap." for
		// Japan.
		"negro", "negroes",
		"jap", "japs",
		"paki", "pakis",
		"honky", "honkey", "honkies",
		"whitey", "whities",
		"gringo", "gringos",
	}},
	{SeverityMild, []string{
		"cracker", "crackers",
		"slope", "slopes",
	}},

	// Homophobic slurs
	{SeveritySevere, []string{
		"fag", "fags", "faggot", "faggots", "faggy",
		"dyke", "dykes", "dykey",
		"poofter", "poofters",
		"buttboy", "buttboys",
		"pillow biter", "pillowbiter",
		"fudgepacker", "fudge packer", "fudgepackers",
		"shirtlifter", "shirt lifter",
		"ass bandit",
	}},
	{SeverityModerate, []string{
		// "Queer" is also a reclaimed self-description; "homo sapiens",
		// Lesbos the island and "sodomite" in scripture are ordinary uses.
		"queer", "queers",
		"homo", "homos",
		"lesbo", "lesbos", "lezbo", "lezbos",
		"sodomite", "sodomites",
		"poof", "poofs",
	}},
	{SeverityMild, []string{
		"fairy", "fairies",
		"pansy", "pansies",
		"queen", "queens",
	}},

	// Transphobic slurs
	{SeveritySevere, []string{
		"tranny", "trannies",
		"heshe", "he-she",
	}},
	{SeverityModerate, []string{
		"shemale", "shemales", "she-male",
		"ladyboy", "ladyboys", "lady-boy",
		"tgirl",
	}},
	{SeverityMild, []string{
		"shim",
		"trap", "traps",
	}},

	// Ableist slurs
	{SeveritySevere, []string{
		"retard", "retards", "retarded", "tard", "tards",
	}},
	{SeverityModerate, []string{
		"spaz", "spazz", "spazzy",
		"spastic",
		"mong", "mongoloid",
	}},
	{SeverityMild, []string{
		"cripple", "crippled", "cripples",
		"mongo",
		"idiot", "idiots", "idiotic",
		"moron", "morons", "moronic",
		"imbecile", "imbeciles",
		"cretin", "cretins",
		"lame", "lamer",
		"dumb", "dummy", "dummies",
		"psycho", "psychos",
		"lunatic", "lunatics",
		"insane",
		"schizo", "schizos",
		"nutjob", "nutjobs",
		"nutcase", "nutcases",
		"freak", "freaks", "freaky",
		"weirdo", "weirdos",
	}},

	// Body-shaming and other crude terms
	{SeverityModerate, []string{
		"midget", "midgets",
		"gimp", "gimps", "gimpy",
	}},
	{SeverityMild, []string{
		"fatso", "fatsos", "fattie", "fatties", "fatty",
		"tubby", "lardo", "porker", "porky",
		"uggo", "uglies", "fugly",
		"dwarf", "dwarfs", "dwarves",
		"runt", "runts",
	}},

	// Anatomical crude terms
	{SeveritySevere, []string{
		"ballsdeep",
		"deepthroat", "deep throat",
	}},
	{SeverityModerate, []string{
		"arsehole", "arseholes", "arse", "arses",
		"ballsack", "ballsacks", "nutsack", "nutsacks",
		"gooch", "grundle",
		"butthole", "buttholes", "buttcrack",
		"clit", "clits",
	}},
	{SeverityMild, []string{
		"anus", "anuses",
		"balls", "nuts",
		"testicle", "testicles", "teste", "testes",
		"scrotum", "scrotums",
		"taint",
		"buttcheek", "buttcheeks",
		"rectum", "rectums",
		"foreskin",
		"clitoris",
		"labia",
		"pube", "pubes", "pubic",
		"pubis",
		"perineum",
	}},

	// Compound and miscellaneous profanity
	{SeveritySevere, []string{
		"milf", "milfs", "dilf", "dilfs",
		"gilf", "gilfs",
		"pawg",
		"thot", "thots",
		"rimjob", "rim job",
		"feltch", "feltching", "felch", "felching",
		"cuntface",
		"thundercunt",
		"pikey", "pikeys",
		"gypo", "gypos", "gyppo", "gyppos",
	}},
	{SeverityModerate, []string{
		"fuckoff", "fuck off", "fuckall", "fuck all",
		"shitload", "shitloads", "shitshow",
		"clusterfuck",
		"mindfuck",
		"clusterf",
		"wtf", "stfu", "gtfo",
		"queef", "queefs", "queefing",
		"douche", "douches", "douchebag", "douchebags", "douchey", "douchy",
		"assmunch", "arsemunch",
		"dickbreath",
		"shitbreath",
		"douchenozzle",
		"jizzface",
		"twatface",
		"dickweasel",
		"wankstain",
		"bellend", "bell end", "bellends",
		"knobhead", "knobheads", "knobend",
		"tosser", "tossers",
		"buggery",
	}},
	{SeverityMild, []string{
		"shart", "sharted", "sharting",
		"scumbag", "scumbags",
		"sleazebag", "sleazeball",
		"dirtbag", "dirtbags",
		"buttmunch", "butthead", "buttwipe",
		"knob", "knobs",
		"minger", "mingers", "minging",
		"git", "gits",
		"sod", "sods", "sodding", "sodoff", "sod off",
		"bollocks", "bollock", "bollox",
		"bloody",
		"bugger", "buggers", "buggered",
		"prat", "prats",
		"pillock", "pillocks",
		"plonker", "plonkers",
		"numpty", "numpties",
		"muppet", "muppets",
		"divvy", "divvies",
		"chav", "chavs", "chavvy",
	}},
}

// profanityList is every listed word, in wordGroups order.
var profanityList = func() []string {
	var words []string
	for _, group := range wordGroups {
		words = append(words, group.words...)
	}
	return words
}()
