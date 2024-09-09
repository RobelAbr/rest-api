package main

// Importiert Standard- und externe Pakete:
// - fmt: für die Formatierung und Ausgabe von Text.
// - github.com/go-resty/resty/v2: Resty ist ein Go HTTP-Client für REST-APIs.
// - log: zum Protokollieren von Fehlern.
import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"log"
)

func main() {
	// Erstellt einen neuen Resty-Client, der für HTTP-Anfragen verwendet wird.
	client := resty.New()

	// Abrufen aller Datensätze und Ausgabe des Ergebnisses.
	fmt.Println("Abrufen aller Datensätze:")
	getAllData(client)

	// Abrufen eines einzelnen Datensatzes mit der ID 1 und Ausgabe des Ergebnisses.
	fmt.Println("\nAbrufen des Datensatzes mit ID 1:")
	getDataByID(client, 1)

	// Abrufen eines einzelnen Datensatzes mit der ID 2 und Ausgabe des Ergebnisses.
	fmt.Println("\nAbrufen des Datensatzes mit ID 2:")
	getDataByID(client, 2)

	// Abrufen eines Datensatzes mit einer ungültigen ID (z. B. 999) und Ausgabe des Ergebnisses.
	fmt.Println("\nAbrufen eines Datensatzes mit ungültiger ID (z.B. 999):")
	getDataByID(client, 999)
}

// getAllData sendet eine GET-Anfrage, um alle Datensätze vom Server abzurufen.
func getAllData(client *resty.Client) {
	// Sende eine GET-Anfrage an den Server mit einem Authorization-Header.
	resp, err := client.R().
		SetHeader("Authorization", "robel"). // Setzt den Authorization-Header auf "robel".
		Get("http://localhost:8080/data")    // URL des Endpunkts für alle Daten.

	// Falls ein Fehler auftritt, wird die Anwendung mit einer Fehlermeldung beendet.
	if err != nil {
		log.Fatalf("Error while getting all data: %v", err)
	}

	// Gibt den HTTP-Statuscode und die Antwort (Body) aus.
	fmt.Println("Status Code:", resp.StatusCode())
	fmt.Println("Body:", string(resp.Body())) // Die Antwort wird in einen String konvertiert und ausgegeben.
}

// getDataByID sendet eine GET-Anfrage, um einen spezifischen Datensatz basierend auf der ID abzurufen.
func getDataByID(client *resty.Client, id int) {
	// Formatiert die URL mit der angegebenen ID (z.B. /data/1).
	url := fmt.Sprintf("http://localhost:8080/data/%d", id)

	// Sende eine GET-Anfrage an den Server mit einem Authorization-Header.
	resp, err := client.R().
		SetHeader("Authorization", "Bearer your_generated_token"). // Setzt den Authorization-Header auf "Bearer your_generated_token".
		Get(url)                                                   // URL für den spezifischen Datensatz.

	// Falls ein Fehler auftritt, wird die Anwendung mit einer Fehlermeldung beendet.
	if err != nil {
		log.Fatalf("Error while getting data by ID: %v", err)
	}

	// Gibt den HTTP-Statuscode und die Antwort (Body) aus.
	fmt.Println("Status Code:", resp.StatusCode())
	fmt.Println("Body:", string(resp.Body())) // Die Antwort wird in einen String konvertiert und ausgegeben.
}
