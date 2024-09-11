package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
)

// generateAsciiArt creates ASCII art from the input text using the specified font.
func generateAsciiArt(input, font string) string {
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
					art = standardFont[' '] // Default to space if character not found
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
					art = shadowFont[' '] // Default to space if character not found
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
					art = thinkertoyFont[' '] // Default to space if character not found
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
		asciiArt := generateAsciiArt(input, font)

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
            handleIndex(w,r)
        }
		// If the path matches, continue to the next handler.
		
	})
}


func main() {
	// Setup HTTP handlers
	http.HandleFunc("/ascii-art", handleGenerateAsciiArt)
	//http.HandleFunc("/", handleIndex)

	http.Handle("/", catchAll(http.DefaultServeMux))
    

	// Log server status and start the server
	fmt.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
//405, //403 ,/500
