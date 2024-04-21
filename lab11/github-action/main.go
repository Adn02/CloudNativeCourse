// Launch microservice server- main.go
package main
import (
	"github.com/Adn02/CloudNativeCourse/lab11/github-action/microservice"
	"log"
)
func main() {
	s := microservice.NewServer("", "8000")
	log.Fatal(s.ListenAndServe())
}