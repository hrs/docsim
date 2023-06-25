package corpus

type Config struct {
	BestFirst      bool
	FollowSymlinks bool
	Limit          int
	NoStemming     bool
	NoStoplist     bool
	ShowScores     bool
	Stoplist       *Stoplist
	Verbose        bool
}
