package exercise1

import (
	"fmt"
	"math/rand"
)

func generateDNAstring(length int, characters string) string {

	// Initialize an empty string to store the result
	result := ""

	// Iterate over the length and append random characters from the given characters
	for i := 0; i < length; i++ {
		randomIndex := rand.Intn(len(characters))
		result += string(characters[randomIndex])
	}

	return result
}

func hammingDistanceCalculate(firstStrand string, secondStrand string) int {
	//Different length return -1
	if len(firstStrand) != len(secondStrand) {
		return -1
	}

	hammingDist := 0

	//Calculate different char or neucleobase to increment hamming distance
	for i := 0; i < len(firstStrand); i++ {
		if firstStrand[i] != secondStrand[i] {
			hammingDist++
		}
	}

	return hammingDist
}

func main() {
	nucleobases := "AGCT"

	for i := 0; i < 1000; i++ {
		firstStrand := generateDNAstring(rand.Intn(11)+10, nucleobases)
		secondStrand := generateDNAstring(rand.Intn(11)+10, nucleobases)

		fmt.Println()
		fmt.Println(firstStrand)
		fmt.Println(secondStrand)

		fmt.Println("Hamming distance: ", hammingDistanceCalculate(firstStrand, secondStrand))
	}
}
