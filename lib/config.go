package main

type Config struct {
	BestFirst      bool
	FollowSymlinks bool
	Limit          int
	NoStemming     bool
	NoStoplist     bool
	OmitQuery      bool
	ShowScores     bool
	Verbose        bool
}
