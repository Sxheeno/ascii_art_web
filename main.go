package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
)

// generateAsciiArt creates ASCII art from the input text using the specified font.
func generateAsciiArt(w http.ResponseWriter, input, font string) string {
	if font != "standard" {
		if font == "shadow" {
			font = "shadow"
		} else if font == "thinkertoy" {
			font = "thinkertoy"
		} else {
			log.Printf("Font %s not found, using default", font)
			font = "standard"
		}
	}

	inputLines := strings.Split(input, "\n")
	var result strings.Builder
	if font == "standard" {
		for _, inputLine := range inputLines {
			lines := make([]string, 8)

			for _, char := range inputLine {
				art, ok := standardFont[char]
				if !ok {
					if char == 32 { // ASCII code for space
						art = standardFont[' ']
					} else {
						log.Printf("Character %c is not a valid ASCII art character", char)
						http.Error(w, fmt.Sprintf("400 - Bad Request: Unsupported character '%c'", char), http.StatusBadRequest)
						return ""
					}
				}

				for i := range art {
					lines[i] += art[i]
				}
			}

			for _, line := range lines {
				result.WriteString(line + "\n")
			}
			result.WriteString("\n")
		}
	} else if font == "shadow" {
		for _, inputLine := range inputLines {
			lines := make([]string, 8)

			for _, char := range inputLine {
				art, ok := shadowFont[char]
				if !ok {
					if char == 32 { // ASCII code for space
						art = standardFont[' ']
					} else {
						log.Printf("Character %c is not a valid ASCII art character", char)
						http.Error(w, fmt.Sprintf("400 - Bad Request: Unsupported character '%c'", char), http.StatusBadRequest)
						return ""
					}
				}

				for i := range art {
					lines[i] += art[i]
				}
			}

			for _, line := range lines {
				result.WriteString(line + "\n")
			}
			result.WriteString("\n")
		}
	} else if font == "thinkertoy" {
		for _, inputLine := range inputLines {
			lines := make([]string, 8)

			for _, char := range inputLine {
				art, ok := thinkertoyFont[char]
				if !ok {
					if char == 32 { // ASCII code for space
						art = standardFont[' ']
					} else {
						log.Printf("Character %c is not a valid ASCII art character", char)
						http.Error(w, fmt.Sprintf("400 - Bad Request: Unsupported character '%c'", char), http.StatusBadRequest)
						return ""
					}
				}

				for i := range art {
					lines[i] += art[i]
				}
			}

			for _, line := range lines {
				result.WriteString(line + "\n")
			}
			result.WriteString("\n")
		}
	}
	return result.String()
}

// Template for displaying ASCII art
var artTemplate = template.Must(template.ParseFiles("static/result.html"))

func handleGenerateAsciiArt(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		input := r.FormValue("input")
		font := r.FormValue("font")
		if font == "" {
			font = "standard"
		}
		log.Printf("Generating ASCII art for input: %s with font: %s", input, font)
		asciiArt := generateAsciiArt(w, input, font)

		if asciiArt == "" {
			return // If an error was already sent, stop further processing
		}

		// Serve the ASCII art using a separate HTML template
		data := struct {
			Art string
		}{
			Art: asciiArt,
		}

		// Render the result page with the ASCII art
		if err := artTemplate.Execute(w, data); err != nil {
			http.Error(w, "Unable to render page", http.StatusInternalServerError)
		}
	} else {
		http.Error(w, "405 Invalid request method", http.StatusMethodNotAllowed)
	}
}

// handleIndex serves the index.html file.
func handleIndex(w http.ResponseWriter, r *http.Request) {
	log.Println("Serving index.html")
	http.ServeFile(w, r, "static/index.html")
}

// handleNotFound handles 404 errors.
func handleNotFound(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "404 - Page Not Found", http.StatusNotFound)
}

// handleError handles other server errors.
func handleError(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "500 - Internal Server Error", http.StatusInternalServerError)
}

// Middleware to catch unknown routes and redirect to 404 handler.
func catchAll(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if the path is exactly "/" or "/ascii-art".
		if r.URL.Path != "/" && r.URL.Path != "/ascii-art" {
			handleNotFound(w, r)
			return
		} else {
			handleIndex(w, r)
		}
	})
}

func loadFont(filename string) error {
	_, err := os.Open(filename)
	if err != nil {
		return err
	}
	return nil
}


func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var err error

		// Load fonts
		err = loadFont("standard.txt")
		if err != nil {
			log.Printf("500 - Internal Server Error: Failed to load standard font: %v", err)
			http.Error(w, "500 - Internal Server Error: Failed to load standard font", http.StatusInternalServerError)
			return
		}

		err = loadFont("shadow.txt")
		if err != nil {
			log.Printf("500 - Internal Server Error: Failed to load shadow font: %v", err)
			http.Error(w, "500 - Internal Server Error: Failed to load shadow font", http.StatusInternalServerError)
			return
		}

		err = loadFont("thinkertoy.txt")
		if err != nil {
			log.Printf("500 - Internal Server Error: Failed to load thinkertoy font: %v", err)
			http.Error(w, "500 - Internal Server Error: Failed to load thinkertoy font", http.StatusInternalServerError)
			return
		}

		catchAll(http.DefaultServeMux).ServeHTTP(w, r)
	})

	http.HandleFunc("/ascii-art", handleGenerateAsciiArt)
	//http.Handle("/", catchAll(http.DefaultServeMux))

	fmt.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

