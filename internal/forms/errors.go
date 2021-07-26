package forms

type errors map[string][]string // create a type that is a map with index of type string and value is slice of strings

// Add adds an error message for a given form field
func (e errors) Add(field, message string) { // receiver is type error, function is add, get field from and also message
	e[field] = append(e[field], message) // in the map, get the list of strings for index field name and append message

}

// Get returns the first error message
func (e errors) Get(field string) string {
	es := e[field]
	if len(es) == 0 {
		return ""
	} else {
		return es[0]
	}
}