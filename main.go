package main

import (
	"database/sql"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

func send(socket io.Writer, data []byte) {
	length := make([]byte, 4)
	binary.LittleEndian.PutUint32(length, uint32(len(data)))

	data = append(length, data...)
	socket.Write(data)
}

func recv(socket io.Reader) []byte {
	buffer := make([]byte, 4)
	socket.Read(buffer)
	length := binary.LittleEndian.Uint32(buffer)

	buffer = make([]byte, length)
	socket.Read(buffer)
	return buffer
}

func contactServerJSON(message map[string]string) (map[string]string, error) {
	remote_addr := new(net.TCPAddr)
	remote_addr.IP = []byte{127, 0, 0, 1}
	remote_addr.Port = 50508

	connection, err := net.DialTCP("tcp4", nil, remote_addr)
	if err != nil {
		return nil, err
	}

	json_message, err := json.Marshal(message)
	if err != nil {
		return nil, err
	}

	send(connection, json_message)
	json_response := recv(connection)

	response := *new(map[string]string)
	json.Unmarshal(json_response, &response)

	return response, nil
}

func generateSignature(payload string) (signature, public_key string) {
	req := map[string]string{"command": "generate", "payload": payload}
	response, _ := contactServerJSON(req)

	return response["signature"], response["public key"]
}

func verifySignature(payload, signature, public_key string) bool {
	req := map[string]string{"command": "verify", "payload": payload, "signature": signature, "public key": public_key}
	response, _ := contactServerJSON(req)

	return response["is valid"] == "True"
}

func blacklistSignature(payload, signature, public_key string) bool {
	req := map[string]string{"command": "blacklist", "payload": payload, "signature": signature, "public key": public_key}
	response, _ := contactServerJSON(req)

	return response["success"] == "True"
}

// ---------------------TODO: Separate These Functions Into Different Files-----------------------------------

type Pages struct {
	db            *sql.DB
	template_path string
}

func pageNotFound(w http.ResponseWriter) {
	w.WriteHeader(404)
	w.Write([]byte("404 - Page Not Found"))
}

func (p *Pages) home(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Called home")

	if r.URL.Path != "/" {
		pageNotFound(w)
		return
	}

	document, _ := template.ParseFiles(p.template_path+"base.html", p.template_path+"home.html")
	document.Execute(w, nil)
}

//obviously for testing only
func (p *Pages) login(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Called login")

	if r.Method == "POST" {
		stmt, err := p.db.Prepare("SELECT Password FROM user_data WHERE Username = ?")
		fmt.Println(err)
		r.ParseForm()

		fmt.Println(r.PostForm["username"][0])

		password := new(string)
		err = stmt.QueryRow(r.PostForm["username"][0]).Scan(password)
		fmt.Println(err)

		if *password == r.PostForm["password"][0] {
			fmt.Println("Authenticated")
		} else {
			fmt.Println("Not Authenticated")
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	document, _ := template.ParseFiles(p.template_path+"base.html", p.template_path+"login.html")
	document.Execute(w, nil)
}

func (p *Pages) signup(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Called signup")

	if r.Method == "POST" {
		stmt, _ := p.db.Prepare("INSERT INTO user_data (Username, Password) VALUES (?, ?)")
		r.ParseForm()

		stmt.Exec(r.PostForm["username"][0], r.PostForm["password"][0])
		defer stmt.Close()

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	document, _ := template.ParseFiles(p.template_path+"base.html", p.template_path+"signup.html")
	document.Execute(w, nil)
}

func main() {
	pages := new(Pages)
	pages.db, _ = sql.Open("mysql", "matthew:MysqlPassword111@tcp(127.0.0.1:3306)/UKIW")
	pages.template_path = "templates/"

	//testng only
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", pages.home)
	http.HandleFunc("/signup", pages.signup)
	http.HandleFunc("/login", pages.login)

	http.ListenAndServe("192.168.1.105:8000", nil)
}
