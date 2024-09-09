package main

// Import von Standardpaketen:
// - encoding/json: zum Arbeiten mit JSON-Daten.
// - fmt: zum Formatieren von Ausgaben.
// - ioutil: zum Lesen von Dateien.
// - net/http: zum Arbeiten mit HTTP-Servern und -Anfragen.
// - strconv: zum Konvertieren von Strings in andere Datentypen (hier int).
import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

// Data repräsentiert die JSON-Datenstruktur.
// Die `json` Tags geben an, wie die JSON-Felder im Go-Struct zugeordnet werden sollen.
type Data struct {
	ID         int    `json:"id"`              // Eindeutige Identifikationsnummer.
	Image      string `json:"jpg"`             // Bild als URL oder Pfad.
	Name       string `json:"name"`            // Nachname der Person.
	FirstName  string `json:"vorname"`         // Vorname der Person.
	Address    string `json:"adresse"`         // Adresse der Person.
	PanCard    string `json:"pan_card_number"` // PAN-Kartennummer der Person.
	ExpiryDate string `json:"expiration_date"` // Ablaufdatum der PAN-Karte.
}

// AuthMiddleware validiert den bereitgestellten Token aus dem Authorization-Header.
// Wenn der Token nicht korrekt ist, wird der Zugriff verweigert.
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extrahiere den Token aus dem Header der HTTP-Anfrage.
		token := r.Header.Get("Authorization")

		// Wenn der Token nicht "robel" ist, verweigere den Zugriff (403 Forbidden).
		if token != "robel" {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		// Wenn der Token korrekt ist, wird die Anfrage an den nächsten Handler weitergeleitet.
		next.ServeHTTP(w, r)
	})
}

func main() {
	// Route für /data, die über die AuthMiddleware geschützt wird und auf handleDataRequest zeigt.
	http.Handle("/data", AuthMiddleware(http.HandlerFunc(handleDataRequest)))

	// Route für /data/{id}, die ebenfalls durch die AuthMiddleware geschützt wird und auf handleDataByIDRequest zeigt.
	http.Handle("/data/{id}", AuthMiddleware(http.HandlerFunc(handleDataByIDRequest)))

	// Ausgabe, dass der Server läuft.
	fmt.Println("Server is running on port 8080")

	// Start des HTTP-Servers auf Port 8080. Nil bedeutet, dass der Default-Mux verwendet wird.
	http.ListenAndServe(":8080", nil)
}

// handleDataRequest bearbeitet Anfragen an /data.
// Es liest die gesamte JSON-Datei und gibt sie als JSON-Antwort zurück.
func handleDataRequest(w http.ResponseWriter, r *http.Request) {
	// Ruft die Daten aus der JSON-Datei ab.
	data, err := readAllJSONData("user.json")
	if err != nil {
		// Falls ein Fehler auftritt, wird ein interner Serverfehler (500) zurückgegeben.
		http.Error(w, "Failed to read data", http.StatusInternalServerError)
		return
	}

	// Setzt den Header auf "application/json", um anzugeben, dass JSON-Daten zurückgegeben werden.
	w.Header().Set("Content-Type", "application/json")

	// Kodiert die Daten als JSON und sendet sie in der HTTP-Antwort.
	json.NewEncoder(w).Encode(data)
}

// handleDataByIDRequest bearbeitet Anfragen an /data/{id}.
// Es extrahiert die ID aus der URL, sucht die entsprechende Person in der JSON-Datei und gibt diese zurück.
func handleDataByIDRequest(w http.ResponseWriter, r *http.Request) {
	// Extrahiere die ID als String aus der URL (nach "/data/").
	idStr := r.URL.Path[len("/data/"):]

	// Konvertiere die ID von einem String in eine Ganzzahl.
	id, err := strconv.Atoi(idStr)
	if err != nil {
		// Wenn die ID keine gültige Zahl ist, wird ein Fehler (400 Bad Request) zurückgegeben.
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	// Liest die gesamte JSON-Datei.
	data, err := readAllJSONData("user.json")
	if err != nil {
		// Falls ein Fehler auftritt, wird ein interner Serverfehler (500) zurückgegeben.
		http.Error(w, "Failed to read data", http.StatusInternalServerError)
		return
	}

	// Durchsuche die Daten nach einem Datensatz mit der gesuchten ID.
	for _, record := range data {
		if record.ID == id {
			// Wenn die ID übereinstimmt, wird der Datensatz als JSON-Antwort zurückgegeben.
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(record)
			return
		}
	}

	// Falls keine Daten für die ID gefunden werden, wird ein Fehler (404 Not Found) zurückgegeben.
	http.Error(w, "Data not found", http.StatusNotFound)
}

// readAllJSONData liest die JSON-Daten aus einer Datei und gibt sie als Slice von Data-Objekten zurück.
func readAllJSONData(filename string) ([]Data, error) {
	// Deklariert eine Variable, um die JSON-Daten zu speichern.
	var data []Data

	// Liest den gesamten Inhalt der Datei in ein Byte-Array ein.
	fileBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		// Wenn ein Fehler beim Lesen der Datei auftritt, wird er zurückgegeben.
		return nil, err
	}

	// Unmarshalled die Byte-Daten in die Data-Struktur.
	err = json.Unmarshal(fileBytes, &data)

	// Gibt die Daten (und ggf. einen Fehler) zurück.
	return data, err
}
