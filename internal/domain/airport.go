package domain

type Location struct {
	En string
	Ru string
}

type Airport struct {
	Code        string
	Name        Location
	City        Location
	Country     Location
	Coordinates [2]float64
	Timezone    string
}
