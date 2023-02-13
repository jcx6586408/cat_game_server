package msg

type Rank struct {
	UID          string
	Val          int
	WorldRank    int64
	CityRank     int64
	Icon         string
	Skin         int
	NickName     string
	Country      string
	CountryShort string
	Province     string
	City         string
	Addr         string
}

type RankSelfRequest struct {
	UID      string
	Val      int
	Icon     string
	NickName string
	Skin     int
}

type BackRankInfo struct {
	WorldRank []*Rank
}

type CityRanks struct {
	Rank []*Rank
}

type RankPull struct {
}
