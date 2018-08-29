package conform

// If returns an Updater that executes one of the given Updater instances, based on the result of a predicate
func If(pred Predicate, t Updater, f Updater) Updater {
	return func(data interface{}) error {
		if pred(data) {
			return t.Do(data)
		}
		return f.Do(data)
	}
}

// IfKey returns an Updater that executes one of the given KeyUpdater instances, based on the result of a predicate
func IfKey(key string, pred KeyPredicate, t KeyUpdater, f KeyUpdater) Updater {
	return func(data interface{}) error {
		if pred(key)(data) {
			return t.Do(key).Do(data)
		}
		return f.Do(key).Do(data)
	}
}

// Compose returns an Updater that executes the given Updater instances in turn.
func Compose(us ...Updater) Updater {
	var all Updater
	for _, u := range us {
		all = all.Then(u)
	}
	return all
}
