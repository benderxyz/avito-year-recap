package aggregation

type Metrics struct {
	ViewsRealty int `json:"views_realty"`
	ViewsAuto   int `json:"views_auto"`
	Purchases   int `json:"purchases"`
	Favorites   int `json:"favorites"`
}

var metricsByUser = map[string]Metrics{
	"1": {ViewsRealty: 87, ViewsAuto: 12, Purchases: 1, Favorites: 15},
	"2": {ViewsRealty: 5, ViewsAuto: 45, Purchases: 6, Favorites: 8},
	"3": {ViewsRealty: 3, ViewsAuto: 2, Purchases: 0, Favorites: 1},
}

func GetByUserID(id string) (Metrics, bool) {
	metrics, ok := metricsByUser[id]
	return metrics, ok
}
