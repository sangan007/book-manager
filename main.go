package main
import(
    "encoding/json"
    "log"
    "net/http"
    "github.com/gorilla/mux"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)
type Book struct {
    ID     uint   `json:"id" gorm:"primaryKey"`
    Title  string `json:"title"`
    Author string `json:"author"`
    Year   int    `json:"year"`
}
var db *gorm.DB
func initDB(){
    dsn := "host=localhost user=postgres password=sangan007 dbname=books_db port=5432 sslmode=disable"
    var err error
    db, err = gorm.Open(postgres.Open(dsn),&gorm.Config{})
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }
    db.AutoMigrate(&Book{})
}
func GetBooks(w http.ResponseWriter, r *http.Request){
    var books []Book
    db.Find(&books)
    json.NewEncoder(w).Encode(books)
}
func GetBookByID(w http.ResponseWriter, r *http.Request){
    params := mux.Vars(r)
    var book Book
    if err := db.First(&book, params["id"]).Error; err != nil {
        http.Error(w, "Book not found", http.StatusNotFound)
        return
    }
    json.NewEncoder(w).Encode(book)
}
func CreateBook(w http.ResponseWriter, r *http.Request){
    var book Book
    if err := json.NewDecoder(r.Body).Decode(&book); err != nil{
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    db.Create(&book)
    json.NewEncoder(w).Encode(book)
}
func UpdateBook(w http.ResponseWriter, r *http.Request){
	params := mux.Vars(r)
    var book Book
    if err := db.First(&book, params["id"]).Error; err != nil{
        http.Error(w, "Book not found", http.StatusNotFound)
        return
    }
    if err := json.NewDecoder(r.Body).Decode(&book); err != nil{
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    db.Save(&book)
    json.NewEncoder(w).Encode(book)
}
func DeleteBook(w http.ResponseWriter, r *http.Request){
    params := mux.Vars(r)
    var book Book
    if err := db.Delete(&book, params["id"]).Error; err != nil{
        http.Error(w, "Book not found", http.StatusNotFound)
        return
    }
    json.NewEncoder(w).Encode("Book deleted successfully")
}
func main(){
    initDB()

    r := mux.NewRouter()
    r.HandleFunc("/books", GetBooks).Methods("GET")
    r.HandleFunc("/books/{id}", GetBookByID).Methods("GET")
    r.HandleFunc("/books", CreateBook).Methods("POST")
    r.HandleFunc("/books/{id}", UpdateBook).Methods("PUT")
    r.HandleFunc("/books/{id}", DeleteBook).Methods("DELETE")

    log.Println("Server started on :8080")
    log.Fatal(http.ListenAndServe(":8080", r))
}
