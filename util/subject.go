package util

var Subjects = map[string]string{
	"Math":       "Math",
	"Biology":    "Biology",
	"Literature": "Literature",
	"Science":    "Science",
	"Physics":    "Physics",
	"Chemical":   "Chemical",
}

func IsSupportedSubject(subject string) bool {
	if _, ok := Subjects[subject]; ok {
		return true
	}
	return false
}
