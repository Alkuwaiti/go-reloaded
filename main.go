package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

// this is a function to remove an element from a slice based on index
func remove(slice []string, s int) []string {
    result := make([]string, len(slice)-1)
    copy(result, slice[:s])
    copy(result[s:], slice[s+1:])
    return result
}

func removeInt(slice []int, s int) []int {
    if s < 0 || s >= len(slice) {
        return slice
    }

    result := make([]int, len(slice)-1)
    copy(result, slice[:s])
    copy(result[s:], slice[s+1:])
    return result
}

// this function lets you insert at a certain index
func insertAtIndex(slice []string, index int, value string) []string {
	// checking if the index is within bounds
	if index < 0 || index > len(slice) {
		return slice
	}

	// Create a new slice with enough capacity to accommodate the new element
	result := make([]string, len(slice)+1)

	// copying the elements before the insertion point
	copy(result[:index], slice[:index])

	// inserting the new element at the specified index
	result[index] = value

	// copying the elements after the insertion point
	copy(result[index+1:], slice[index:])

	return result
}

func isNumberWithClosingParenthesis(input string) bool {
	// defining a regular expression pattern to match a number followed by a closing parenthesis
	re := regexp.MustCompile(`^\d+\)$`)

	// check if the input string matches the pattern which returns true or false
	return re.MatchString(input)
}

func extractNumber(input string) string {
	// defining a regular expression pattern to match digits
	re := regexp.MustCompile(`\d+`)

	// finding all matches in the input string
	matches := re.FindStringSubmatch(input)

	// checking if there is at least one match
	if len(matches) > 0 {
		return matches[0] // the first match is the numeric part
	}

	// No match found, return an empty string or handle the case as needed
	return ""
}


func defaultProcessing(body *string, parameter string) {
	// .Fields splits a string into a slice of strings
	sliceOfStrings := strings.Fields(*body)
	
	// saving it in a new variable to also avoid indexing issues
	newSliceOfStrings := sliceOfStrings
	
	baseInNumber := 16
	if parameter == "(bin)" {
		baseInNumber = 2
	}
	
	// loop to find every word in the sliceOfStrings
	// starting from the right to avoid indexing issues by removing the elements.
	for i := len(sliceOfStrings)-1; i>=0; i-- {
		
		// check if a word is (hex) or (bin)
		if sliceOfStrings[i] == parameter {

			// check if the index is 0
			if i==0 {
				fmt.Println("Parameter " + parameter + " can't be at the start of a string")
				os.Exit(1);
			}
			
			// get the actual number before (hex) and save it in a variable
			hexString := sliceOfStrings[i-1]
			
			// remove the word (hex), and save the slice in a new var to avoid indexing issues
			newSliceOfStrings = remove(newSliceOfStrings, i)
			
			// extract the decimal with parseInt
			decimalInt, err := strconv.ParseInt(hexString, baseInNumber, 64)
			
			// error handling for parseInt
			if err != nil {
				fmt.Println("Error during converting the ParseInt")
				fmt.Println(i)
				return
			}
			
			// remove the hex number from the slice
			newSliceOfStrings = remove(newSliceOfStrings, i-1)
			
			// convert the decimalInt to a string to add again to the slice later on
			decimalString := strconv.FormatInt(decimalInt, 10)
			
			// insert the decimal value at the specified index
			newSliceOfStrings = insertAtIndex(newSliceOfStrings, i-1, decimalString)
			
			// save the string at where the body's address is
			*body = strings.Join(newSliceOfStrings, " ")
			
		}
	}
}

func multipleChanges(body *string, parameter string, i int, newSliceOfStrings []string, numberInInt64 int64) {

		// a loop for the number of integers.
		for j := 1; j <= int(numberInInt64); j++ {
			
			// get the string with i-j
			currentString := newSliceOfStrings[i-j]

			// check if the parameter is (low,
			if parameter == "(low," {
				
				// turn the current string to lowercase
				currentString = strings.ToLower(currentString)

			// check if the parameter is (up,
			} else if parameter == "(up," {
				
				// turn the current string to uppercase
				currentString = strings.ToUpper(currentString)

			// check if the parameter is (cap,	
			} else if parameter == "(cap," {
				// turn the current string to lowercase
				currentString = strings.ToLower(currentString)
				currentString = strings.Title(currentString)
			}

			// remove the currentString's older version
			newSliceOfStrings = remove(newSliceOfStrings, i-j)

			// insert the currentString's new version
			newSliceOfStrings = insertAtIndex(newSliceOfStrings, i-j, currentString)
			
			// save everything in the body
			*body = strings.Join(newSliceOfStrings, " ")
			
		
	}
}

func changeStuff(body *string, parameter string) {
	// .Fields splits a string into a slice of strings
	sliceOfStrings := strings.Fields(*body)
	
	// saving it in a new variable to also avoid indexing issues
	newSliceOfStrings := sliceOfStrings

	// fmt.Println(newSliceOfStrings)
	// fmt.Println(parameter)

	// loop to find every word in the sliceOfStrings
	// starting from the right to avoid indexing issues by removing the elements.
	for i := len(sliceOfStrings)-1; i>=0; i-- {

		// check if a word is [(low, or (up, or (cap,] or (low) or (up) or (cap) making it not go to the else statements
		if sliceOfStrings[i] == parameter {

			// check if the second (next) part is a number with a closing parenthesis
			if len(sliceOfStrings) > i + 1 && isNumberWithClosingParenthesis(sliceOfStrings[i+1]) == true {

				// extract just the number, still a string ...
				numberInString := extractNumber(newSliceOfStrings[i+1])

				// parse the int64 from the string
				numberInInt64, _ := strconv.ParseInt(numberInString, 10, 64)

				// if number of strings from beginning to right before (low, is less than numberInInt64
				if (len(sliceOfStrings[:i]) < int(numberInInt64)) {
					fmt.Println("Inserted number is too high! Number: " + numberInString)
					os.Exit(1)
				}

				// remove the number with parenthesis
				newSliceOfStrings = remove(newSliceOfStrings, i+1)
				
				// remove the parameter (low,
				newSliceOfStrings = remove(newSliceOfStrings, i)
				
				// if parameter is (low,
				if parameter == "(low," {
					
					// call the function with the "(low," parameter
					multipleChanges(body, "(low,", i, newSliceOfStrings, numberInInt64)
				
				// if parameter is (up,
				} else if parameter == "(up," {
					
					// call the function with the "(up" parameter
					multipleChanges(body, "(up,", i, newSliceOfStrings, numberInInt64)

				} else if parameter == "(cap," {
					// call the function with the "(cap," parameter
					multipleChanges(body, "(cap,", i, newSliceOfStrings, numberInInt64)
				}

			}
			// TODO: cap isn't done
		// check if the current element is (up) or (low) or (cap)
		} 
		if sliceOfStrings[i] == "(up)" || sliceOfStrings[i] == "(low)" || sliceOfStrings[i] == "(cap)"{
			
			// get the string to be changed
			stringToBeChanged := newSliceOfStrings[i-1]

			// remove the parameter itself like (low) (up) (cap)
			newSliceOfStrings = remove(newSliceOfStrings, i)

			// remove the stringToBeChanged's old version
			newSliceOfStrings = remove(newSliceOfStrings, i-1)

			// if it is up
			if sliceOfStrings[i] == "(up)" {

				// turn to upper
				stringToBeChanged = strings.ToUpper(stringToBeChanged)
			
			// if it is low
			} else if sliceOfStrings[i] == "(low)" {

				// turn to lower
				stringToBeChanged = strings.ToLower(stringToBeChanged)

			// if it is cap
			} else if sliceOfStrings[i] == "(cap)" {
				// turn to lower
				stringToBeChanged = strings.ToLower(stringToBeChanged)
				// turn to title cap
				stringToBeChanged = strings.Title(stringToBeChanged)
			}

			// insert the stringToBeChanged's new version
			newSliceOfStrings = insertAtIndex(newSliceOfStrings, i-1, stringToBeChanged)

			// save to body
			*body = strings.Join(newSliceOfStrings, " ")
			
		}

	}
}

func hex(body *string)  {
	defaultProcessing(body, "(hex)")
}

func bin(body *string)  {
	defaultProcessing(body, "(bin)")
}

func low(body *string) {
	changeStuff(body, "(low)")
	changeStuff(body, "(low,")

}

func up(body *string) {
	changeStuff(body, "(up)")
	changeStuff(body, "(up,")
}

func cap(body *string) {
	changeStuff(body, "(cap,")
	changeStuff(body, "(cap)")
}

func punctuation(body *string) {
	// .Fields splits a string into a slice of strings
	sliceOfStrings := strings.Fields(*body)
	
	// saving it in a new variable to also avoid indexing issues
	newSliceOfStrings := sliceOfStrings

	// loop to find every word in the sliceOfStrings
	for i := len(sliceOfStrings)-1; i>=0; i-- {

		// check if an indice is . , ! ? : or ; ,
		if sliceOfStrings[i] == "." || sliceOfStrings[i] == "," || sliceOfStrings[i] == "!" || sliceOfStrings[i] == "?" || sliceOfStrings[i] == ":" || sliceOfStrings[i] == ";" || sliceOfStrings[i] == "..." || sliceOfStrings[i] == "!?" {
			

			// get the word to the left of punc mark
			wordToBeModified := newSliceOfStrings[i-1]

			// modify the word by concatination of word + punc mark
			wordToBeModified = wordToBeModified + newSliceOfStrings[i]

			// remove the original punc mark
			newSliceOfStrings = remove(newSliceOfStrings, i)

			// remove the old word to replace it later on
			newSliceOfStrings = remove(newSliceOfStrings, i-1)

			newSliceOfStrings = insertAtIndex(newSliceOfStrings, i-1, wordToBeModified)
		
		// else check if first indice of an indice is . , ! : or ;
		} else if sliceOfStrings[i][0] == '.' || sliceOfStrings[i][0] == ',' || sliceOfStrings[i][0] == '!' || sliceOfStrings[i][0] == ':' || sliceOfStrings[i][0] == ';' {
			
			// get the punctuation mark
			punctuationMark := string(sliceOfStrings[i][0])

			// get the word to be modified which is of the same index
			wordToBeModified := newSliceOfStrings[i]

			// get the word that needs to have punc added
			wordWithNewPunc := newSliceOfStrings[i-1]
			
			// add the punc to the wordWithNewPunc
			wordWithNewPunc = wordWithNewPunc + string(sliceOfStrings[i][0])

			// turn the word to be modified to a slice to make it easier to work with
			sliceOfWord := strings.Split(wordToBeModified, "")

			// remove the first part of the word
			sliceOfWord = remove(sliceOfWord, 0)

			// check if the punc marks are actually ...
			if (len(sliceOfStrings[i]) >= 3 && sliceOfStrings[i][0] == '.' && sliceOfStrings[i][1] == '.' && sliceOfStrings[i][2] == '.')  {

				// modify the punctuationMark if if statement is true
				punctuationMark = "..."

				// remove the second part of the word
				sliceOfWord = remove(sliceOfWord, 0)

				// remove the third part of the word
				sliceOfWord = remove(sliceOfWord, 0)

				// override the wordWithNewPunc with word to left + punctuation mark
				wordWithNewPunc = newSliceOfStrings[i-1] + punctuationMark

			// else check if the punc marks is actually !?
			} else if (len(sliceOfStrings[i]) >= 2 && sliceOfStrings[i][0] == '!' && sliceOfStrings[i][1] == '?') {

				// modify the punctuationMark if if statement is true
				punctuationMark = "!?"

				// remove the second part of the word
				sliceOfWord = remove(sliceOfWord, 0)

				// override the wordWithNewPunc with word to left + punctuation mark
				wordWithNewPunc = newSliceOfStrings[i-1] + punctuationMark
			}

			// join back the wordToBeModified
			wordToBeModified = strings.Join(sliceOfWord, "")

			// remove the word with the punc mark to its left (,don't)
			newSliceOfStrings = remove(newSliceOfStrings, i)

			// insert the new modified word without the left punc mark (don't)
			newSliceOfStrings = insertAtIndex(newSliceOfStrings, i, wordToBeModified)

			// remove the word with the no punctuation (boring)
			newSliceOfStrings = remove(newSliceOfStrings, i-1)
			
			// insert the word with the new punc (boring,)
			newSliceOfStrings = insertAtIndex(newSliceOfStrings, i-1, wordWithNewPunc)
		}
		
		// save to body
		*body = strings.Join(newSliceOfStrings, " ")
	}
}

func quotations(body *string) {
	//TODO :- a word with single quotes
	// .Fields splits a string into a slice of strings
	sliceOfStrings := strings.Fields(*body)
	
	// saving it in a new variable to also avoid indexing issues
	newSliceOfStrings := sliceOfStrings

	if !slices.Contains(newSliceOfStrings, "'"){
		return
	}

	// a slice to save the indexes of every single quote
	quoteIndexes := make([]int, 0)

	// loop to find every ' (single quote) in the sliceOfStrings from the left
	for i := 0; i<len(sliceOfStrings); i++ {
		if sliceOfStrings[i] == "'" {
			quoteIndexes = append(quoteIndexes, i)
		}
	}


	// Check if there are at least two single quotes
	if len(quoteIndexes) < 2 {
		return
	}

	// to use as index for quoteIndexes
	j := len(quoteIndexes) - 1

	// for as long as there are entries in the quoteIndexes slice
	for len(quoteIndexes) > 0 {

		// get the word that is in the right of the quote
		wordinTheRight := newSliceOfStrings[quoteIndexes[j] - 1]

		// fmt.Println("wordinTheRight")
		// modify word in the right to include a ' at the end of it
		wordinTheRight = wordinTheRight + "'"

		// fmt.Println(wordinTheRight)


		// get the word that is in the left of the quote
		wordinTheLeft := newSliceOfStrings[quoteIndexes[j-1] + 1]

		// modify word in the left to include a ' in the beginning of it
		wordinTheLeft = "'" + wordinTheLeft

		// remove the quote in the right from newSliceOfStrings
		newSliceOfStrings = remove(newSliceOfStrings, quoteIndexes[j])

		// // remove the old version of the word in the right
		newSliceOfStrings = remove(newSliceOfStrings, quoteIndexes[j] - 1)

		// insert the new version of the word in the right with the quote '
		newSliceOfStrings = insertAtIndex(newSliceOfStrings, quoteIndexes[j] - 1, wordinTheRight)

		// // remove the old version of the word in the left
		newSliceOfStrings = remove(newSliceOfStrings, quoteIndexes[j-1] + 1)

		// remove the quote in the left from newSliceOfStrings
		newSliceOfStrings = remove(newSliceOfStrings, quoteIndexes[j-1])

		// insert the new version of the word in the left with the quote '
		newSliceOfStrings = insertAtIndex(newSliceOfStrings, quoteIndexes[j-1], wordinTheLeft)

		// fmt.Println("wordinTheLeft")
		// fmt.Println(wordinTheLeft)


		// remove the two rightmost indexes of quotes
		quoteIndexes = removeInt(quoteIndexes, j)
		quoteIndexes = removeInt(quoteIndexes, j-1)

		// let j go back twice so it's not out of range
		j = j -2
		*body = strings.Join(newSliceOfStrings, " ")


		// fmt.Println(" final")
		// fmt.Println(quoteIndexes)
		
	}


	

	

	

	// TODO :- 'awesome' doesn't work because the thing is either 'awesome OR awesome'
	
}

func vowels(body *string) {
	// .Fields splits a string into a slice of strings
	sliceOfStrings := strings.Fields(*body)
	
	// saving it in a new variable to also avoid indexing issues
	newSliceOfStrings := sliceOfStrings

	// loop to find every word in the sliceOfStrings
	for i := len(sliceOfStrings)-1; i>=0; i-- {
		// check if the [0] index of the current word is a vowel.
		if string(sliceOfStrings[i][0]) == "a" || string(sliceOfStrings[i][0]) == "e" || string(sliceOfStrings[i][0]) == "i" || string(sliceOfStrings[i][0]) == "o" || string(sliceOfStrings[i][0]) == "" {

			// check if the word to the left is an "a" or "A"
			if sliceOfStrings[i-1] == "a" || sliceOfStrings[i-1] == "A" {

				// get the a or A so we can just add an "n" to it later
				aToBeChanged := newSliceOfStrings[i-1]

				// add the "n" to it
				aToBeChanged = aToBeChanged + "n"

				// remove the item a or A
				newSliceOfStrings = remove(newSliceOfStrings, i-1)

				// insert the new aToBeChanged at that index
				newSliceOfStrings = insertAtIndex(newSliceOfStrings, i-1, aToBeChanged)

				// save in the body
				*body = strings.Join(newSliceOfStrings, " ")

			}
		}
	}
}

func main() {

	if len(os.Args) != 3 {
        fmt.Println("Usage: go run . input_file output_file")
        os.Exit(1)
    }

	// read the file based on Args
	bodyInBytes, err := ioutil.ReadFile(os.Args[1])

	// if an error exists, kill it
	if err != nil {
        log.Fatalf("unable to read file: %v", err)
    }
	// stringify the body
	body := string(bodyInBytes)

	// passing in the address to the body so that we don't have to return anything with the function call, calling hex function
	hex(&body)

	// call the bin function in the same way
	bin(&body)

	// call the low function
	low(&body)

	// call the up funcion
	up(&body)

	// call the cap funcion
	cap(&body)

	// call the punctuation funcion
	punctuation(&body)

	// call the quotations function
	quotations(&body)

	vowels(&body)

	// create the file
	f, err := os.Create(os.Args[2])
    if err != nil {
        log.Fatal(err)
    }
    

	// write the content to the file
	f.WriteString(body)
	defer f.Close()

	fmt.Println("-----------------------------------------------------------------------")
	fmt.Println(body)
}