package prob

// assert is a helper function to provide
// moderate runtime type checking on the Element interface
func assert(flag bool, s string) {
	if !flag {
		panic(s)
	}
}
