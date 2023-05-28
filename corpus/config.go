package corpus

type Config struct {
	BestFirst      bool
	FollowSymlinks bool
	Limit          int
	NoStemming     bool
	NoStoplist     bool
	OmitQuery      bool
	ShowScores     bool
	Stoplist       *Stoplist
	Verbose        bool
}
