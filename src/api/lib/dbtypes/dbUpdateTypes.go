package dbtypes

type UpdateCase struct {
	Filter map[string]interface{}
	Update map[string]interface{}
}
