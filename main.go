package main

// Importing the functions
import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

type Book struct {
	ID            int    `json:"id"`
	Title         string `json:"title"`
	Author        string `json:"author"`
	Publisher     string `json:"publisher"`
	YearPublished string `json:"year_published"`
}

type BookManager struct {
	jsonFile string
}

// func (b *Book) getJSON() []byte {
// 	data, err := json.MarshalIndent(b, "", "  ")
// 	if err != nil {
// 		panic(err)
// 	}
// 	return data
// }

func NewBookManager(jsonFile string) *BookManager {
	// If jsonFile does not exist, create an empty file
	if _, err := os.Stat(jsonFile); os.IsNotExist(err) {
		err := ioutil.WriteFile(jsonFile, []byte("[]"), 0644)
		if err != nil {
			panic(err)
		}
	}

	return &BookManager{
		jsonFile: jsonFile,
	}
}

func (bm *BookManager) validYear() string {
	for {
		fmt.Print("Enter year published: ")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		yearPublishedStr := scanner.Text()

		yearPublished, err := strconv.Atoi(yearPublishedStr)
		if err != nil {
			fmt.Println("-------------------------------------------")
			fmt.Println("Invalid input. Please enter a valid published year.")
			fmt.Println("-------------------------------------------")
			continue
		}

		// Validate the year is within the accepted range
		if yearPublished < 1900 || yearPublished > 2023 {
			fmt.Println("-------------------------------------------")
			fmt.Println("Invalid year. Please enter a year between 1900 and 2023.")
			fmt.Println("-------------------------------------------")
			continue
		}

		return yearPublishedStr
	}
}

// Adding Book
func (bm *BookManager) addBook() {
	// Open the JSON file for reading and writing
	file, err := os.OpenFile(bm.jsonFile, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	book := &Book{}

	decoder := json.NewDecoder(file)
	encoder := json.NewEncoder(file)

	var books []Book
	err = decoder.Decode(&books)
	if err != nil && err.Error() != "EOF" {
		panic(err)
	}

	// Prompt user for information about the new book
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Enter title of book: ")
	scanner.Scan()
	book.Title = scanner.Text()

	fmt.Print("Enter author of book: ")
	scanner.Scan()
	book.Author = scanner.Text()

	fmt.Print("Enter publisher of book: ")
	scanner.Scan()
	book.Publisher = scanner.Text()

	book.YearPublished = bm.validYear()
	book.ID = len(books) + 1

	books = append(books, *book)

	file.Seek(0, 0)
	err = encoder.Encode(books)
	if err != nil {
		panic(err)
	}
	fmt.Println("Book added successfully.")
}

// ListBook
func (bm *BookManager) listBooks(pageSize int) {
	// Open the JSON file for reading
	file, err := os.Open(bm.jsonFile)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)

	var books []Book
	err = decoder.Decode(&books)
	if err != nil {
		panic(err)
	}

	numPages := len(books) / pageSize
	if len(books)%pageSize != 0 {
		numPages++
	}

	for i := 0; i < numPages; i++ {
		fmt.Printf("\t\tPage %d\n", i+1)
		fmt.Println("-------------------------------------------")

		start := i * pageSize
		end := (i + 1) * pageSize
		if end > len(books) {
			end = len(books)
		}

		for _, book := range books[start:end] {
			fmt.Printf("%d. %s\n", book.ID, book.Title)
		}

		if i != numPages-1 {
			fmt.Println("-------------------------------------------")
			fmt.Print("Press any key for next page (ex. 'n'): ")
			scanner := bufio.NewScanner(os.Stdin)
			scanner.Scan()
			fmt.Println("-------------------------------------------")
		}
	}
}

// Update the Book Using ID
func (bm *BookManager) updateBook(id int) bool {
	// Open the JSON file for reading and writing
	file, err := os.OpenFile(bm.jsonFile, os.O_RDWR, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	encoder := json.NewEncoder(file)

	var books []Book
	err = decoder.Decode(&books)
	if err != nil {
		panic(err)
	}

	for i, book := range books {
		if book.ID == id {
			fmt.Printf("\n----- Current Book Details -----\nTitle: %s\nAuthor: %s\nPublisher: %s\nPublish Year: %s\n-------------------------------------------\n",
				book.Title, book.Author, book.Publisher, book.YearPublished)

			// Prompt the user to enter the new details
			scanner := bufio.NewScanner(os.Stdin)

			fmt.Print("Enter new title: ")
			scanner.Scan()
			newTitle := scanner.Text()

			fmt.Print("Enter new author name: ")
			scanner.Scan()
			newAuthor := scanner.Text()

			fmt.Print("Enter new publisher name: ")
			scanner.Scan()
			newPublisher := scanner.Text()

			newYearPublished := bm.validYear()

			books[i].Title = newTitle
			books[i].Author = newAuthor
			books[i].Publisher = newPublisher
			books[i].YearPublished = newYearPublished

			file.Seek(0, 0)
			err = encoder.Encode(books)
			if err != nil {
				panic(err)
			}

			return true
		}
	}

	return false
}

// Deleting book using ID
func (bm *BookManager) deleteBook(id int) bool {
	// Open the JSON file for reading and writing
	file, err := os.OpenFile(bm.jsonFile, os.O_RDWR, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	encoder := json.NewEncoder(file)

	var books []Book
	err = decoder.Decode(&books)
	if err != nil {
		panic(err)
	}

	for i, book := range books {
		if book.ID == id {
			books = append(books[:i], books[i+1:]...)

			file.Seek(0, 0)
			err = encoder.Encode(books)
			if err != nil {
				panic(err)
			}

			fmt.Println("-------------------------------------------")
			fmt.Printf("Recently deleted book: %s\n", book.Title)

			return true
		}
	}

	return false
}

// View one by one book details that you want
func (bm *BookManager) viewBook(id int) bool {
	// Open the JSON file for reading
	file, err := os.Open(bm.jsonFile)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)

	var books []Book
	err = decoder.Decode(&books)
	if err != nil {
		panic(err)
	}

	for _, book := range books {
		if book.ID == id {
			fmt.Printf("\n----- Book Details -----\nTitle: %s\nAuthor: %s\nPublisher: %s\nPublish Year: %s\n",
				book.Title, book.Author, book.Publisher, book.YearPublished)
			return true
		}
	}

	return false
}

// List the authors that are present
func (bm *BookManager) listAuthors() {
	// Open the JSON file for reading
	file, err := os.Open(bm.jsonFile)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)

	var books []Book
	err = decoder.Decode(&books)
	if err != nil {
		panic(err)
	}

	authors := make(map[string]bool)
	for _, book := range books {
		authors[book.Author] = true
	}

	for author := range authors {
		fmt.Println(author)
	}
}

// Main Function where program get executing
func main() {
	bm := NewBookManager("books.json")

	for {
		fmt.Println("-------------------------------------------")
		fmt.Println("\n\tBook Management Tool")
		fmt.Println("-------------------------------------------")
		fmt.Println("\t*** Menu ***")
		fmt.Println("-------------------------------------------")
		fmt.Println("1. Add Book")
		fmt.Println("2. List Books")
		fmt.Println("3. Update Book")
		fmt.Println("4. Delete Book")
		fmt.Println("5. List Authors")
		fmt.Println("6. View Book Details")
		fmt.Println("7. Quit")
		fmt.Println("-------------------------------------------")

		fmt.Print("Enter your choice (1 to 7): ")
		reader := bufio.NewReader(os.Stdin)
		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)

		fmt.Println("-------------------------------------------")

		switch choice {
		case "1":
			fmt.Println("Add Book")
			fmt.Println("-------------------------------------------")
			bm.addBook()
			fmt.Println("-------------------------------------------")
			fmt.Println("Book added successfully")

		case "2":
			fmt.Println("List Books")
			fmt.Println("-------------------------------------------")
			bm.listBooks(50)

		case "3":
			fmt.Println("Update Book")
			fmt.Println("-------------------------------------------")
			fmt.Print("Enter ID of book to update: ")
			idString, _ := reader.ReadString('\n')
			idString = strings.TrimSpace(idString)
			id, _ := strconv.Atoi(idString)
			if bm.updateBook(id) {
				fmt.Println("-------------------------------------------")
				fmt.Println("Book updated successfully")
			} else {
				fmt.Println("-------------------------------------------")
				fmt.Println("Book not found")
			}

		case "4":
			fmt.Println("Delete Book")
			fmt.Println("-------------------------------------------")
			fmt.Print("Enter ID of book to delete: ")
			idString, _ := reader.ReadString('\n')
			idString = strings.TrimSpace(idString)
			id, _ := strconv.Atoi(idString)
			if bm.deleteBook(id) {
				fmt.Println("-------------------------------------------")
				fmt.Println("Book deleted successfully")
			} else {
				fmt.Println("-------------------------------------------")
				fmt.Println("Book not found")
			}

		case "5":
			fmt.Println("List Authors")
			fmt.Println("-------------------------------------------")
			bm.listAuthors()

		case "6":
			fmt.Println("View Book Details")
			fmt.Println("-------------------------------------------")
			fmt.Print("Enter ID of book to view details: ")
			idString, _ := reader.ReadString('\n')
			idString = strings.TrimSpace(idString)
			id, _ := strconv.Atoi(idString)
			if bm.viewBook(id) {
				fmt.Println("-------------------------------------------")
				fmt.Println("Book found")
			} else {
				fmt.Println("-------------------------------------------")
				fmt.Println("Book not found")
			}

		case "7":
			fmt.Println("Goodbye!")
			return

		default:
			fmt.Println("Invalid choice. Please enter a number from 1 to 7.")
		}
	}
}
