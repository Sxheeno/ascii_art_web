package main

import (
	"fmt"
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

// handleGenerateAsciiArt handles the ASCII art generation requests.
func handleGenerateAsciiArt(w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodPost {
        input := r.FormValue("input")
        font := r.FormValue("font")
        if font == "" {
            font = "standard"
        }
        log.Printf("Generating ASCII art for input: %s with font: %s", input, font)
        asciiArt := generateAsciiArt(input, font)
        fmt.Print(asciiArt)
        w.Header().Set("Content-Type", "text/plain")
        fmt.Fprintf(w, "%s", asciiArt)
    } else {
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
    }
}

// handleIndex serves the index.html file.
func handleIndex(w http.ResponseWriter, r *http.Request) {
    log.Println("Serving index.html")
    http.ServeFile(w, r, "static/index.html")
}

func main() {
    // Setup HTTP handlers
    http.HandleFunc("/generate-ascii-art", handleGenerateAsciiArt)
    http.HandleFunc("/", handleIndex)

    fmt.Println("Server is running on port 8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
