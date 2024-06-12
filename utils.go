package main

func ToString(strings []string) string {
	result := ""
	for _, str := range strings {
		result += str + "\n"
	}
	return result
}
