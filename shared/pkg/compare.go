package shared

import "fmt"

func Compare(
	laterFile string,
	earlierFile string,
) (
	jsonFileLocation string,
	err error,
) {
	fmt.Printf("Comparing %q to %q\n", laterFile, earlierFile)

	return "", nil
}
