package msg

type Rank struct {
	UID          string
	Val          int
	Icon         string
	NickName     string
	Country      string
	CountryShort string
	City         string
}

type RankSelfRequest struct {
	UID      string
	Val      int
	Icon     string
	NickName string
}

type BackRankInfo struct {
	WorldRank []*Rank
}

type CityRanks struct {
	Rank []*Rank
}

type RankPull struct {
}
