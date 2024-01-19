package mongodb

import "go.mongodb.org/mongo-driver/mongo/options"

type MongoPaginate struct {
	Limit      int64
	Page       int64
	TotalRows  int64
	TotalPages int64
}

func NewMongoPaginate(limit int, page int) *MongoPaginate {
	return &MongoPaginate{
		Limit: int64(limit),
		Page:  int64(page),
	}
}

func (mp *MongoPaginate) GetPaginatedOpts() *options.FindOptions {
	l := mp.Limit
	skip := mp.Page*mp.Limit - mp.Limit
	fOpt := options.FindOptions{Limit: &l, Skip: &skip}

	return &fOpt
}

func (mp *MongoPaginate) CalculateTotalPages(totalRows int64) int64 {
    if mp.Limit == 0 {
        return 0
    }
    mp.TotalRows = totalRows
    mp.TotalPages = (totalRows + mp.Limit - 1) / mp.Limit
    return mp.TotalPages
}
