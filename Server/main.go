package main

import (
	signing "UKIWcoursework/Server/Signing"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

type UserDetails struct {
	username string
}

var errTokenInvalid error = errors.New("the given token is invalid")
var errTokenExpired error = errors.New("the given token is expired")

func parseToken(cookie *http.Cookie) (user_details *UserDetails, err error) {
	unescaped_token, err := url.PathUnescape(cookie.Value)
	if err != nil {
		return nil, err
	}

	json_token := []byte(unescaped_token)

	token := new(Token)
	err = json.Unmarshal(json_token, token)
	if err != nil {
		return nil, err
	}

	expiration, err := strconv.Atoi(token.Expiration)
	if err != nil {
		return nil, err
	}

	if expiration < int(time.Now().Unix()) {
		return nil, errTokenExpired
	}

	payload := map[string]string{
		"username":   token.Username,
		"expiration": token.Expiration,
	}

	json_payload, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	verified, err := signing.VerifySignature(string(json_payload), token.Signature, token.Public_key)
	if err != nil {
		return nil, err
	}

	if !verified {
		return nil, errTokenInvalid
	}

	user_details = &UserDetails{token.Username}
	return user_details, nil
}

type handler func(w http.ResponseWriter, r *http.Request, user_details *UserDetails) ErrorResponse

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var http_error ErrorResponse

	cookie, err := r.Cookie("auth_token")
	//a cookie has been sent
	if err == nil {
		user_details, err := parseToken(cookie)
		fmt.Println(err)
		//token is ok, run function and pass user details
		if err == nil {
			http_error = h(w, r, user_details)

		} else if err == errTokenExpired || err == errTokenInvalid {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return

		} else {
			//some other error has happened
			http_error = HTTPerror{500}
		}
	} else {
		//no cookie has been sent
		http_error = h(w, r, nil)
	}

	if http_error == nil {
		return
	}

	//handle error
	w.Header().Add("content-type", "text/html")
	w.WriteHeader(http_error.Code())

	message := "<h1>" + http_error.Error() + "</h1>"
	w.Write([]byte(message))
}

type ErrorResponse interface {
	Code() int
	Error() string
}

type HTTPerror struct {
	code int
}

func (e HTTPerror) Code() int {
	return e.code
}

func (e HTTPerror) Error() string {
	switch e.code {
	case 404:
		return "404 - Page Not Found"
	case 500:
		return "500 - Internal Server Error"
	}

	return ""
}

type Pages struct {
	db            *sql.DB
	template_path string
}

func (p *Pages) home(w http.ResponseWriter, r *http.Request, user_details *UserDetails) ErrorResponse {
	fmt.Println("Called home")

	if r.URL.Path != "/" {
		return HTTPerror{404}
	}

	document, _ := template.ParseFiles(p.template_path+"base.html", p.template_path+"home.html")
	document.Execute(w, nil)
	return nil
}

type Token struct {
	Username   string
	Expiration string
	Signature  string
	Public_key string
}

func loginUser(w http.ResponseWriter, username string) error {
	//20 days
	expiration := time.Now().Unix() + 28800

	payload := map[string]string{
		"username":   username,
		"expiration": strconv.Itoa(int(expiration)),
	}

	json_payload, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	signature, public_key, err := signing.GenerateSignature(string(json_payload))
	if err != nil {
		return err
	}

	token := Token{
		Username:   username,
		Expiration: payload["expiration"],
		Signature:  signature,
		Public_key: public_key,
	}

	json_token, err := json.Marshal(token)
	if err != nil {
		return err
	}

	cookie := new(http.Cookie)
	cookie.Name = "auth_token"
	cookie.Value = url.PathEscape(string(json_token))

	http.SetCookie(w, cookie)
	return nil
}

type 

//obviously for testing only
func (p *Pages) login(w http.ResponseWriter, r *http.Request, user_details *UserDetails) ErrorResponse {
	fmt.Println("Called login")

	
	if err != nil {
		return HTTPerror{500}
	}

	if r.Method == "POST" {
		stmt, err := p.db.Prepare("SELECT Password FROM user_data WHERE Username = ?")
		if err != nil {
			return HTTPerror{500}
		}

		err = r.ParseForm()
		if err != nil {
			return HTTPerror{500}
		}

		username := r.PostForm["username"][0]
		raw_password := r.PostForm["password"][0]
		database_hash := new(string)

		err = stmt.QueryRow(username).Scan(database_hash)
		if err != nil {
			
			data := LoginTemplateData{
				*user_details,
				"The username you entered does not exist!",
			}
			err = document.Execute(w, data)

			if err != nil {
				return HTTPerror{500}
			}
			return nil
		}

		err = bcrypt.CompareHashAndPassword([]byte(*database_hash), []byte(raw_password))
		if err == nil {
			fmt.Println("Authenticated")

			err = loginUser(w, username)
			if err != nil {
				return HTTPerror{500}
			}

			http.Redirect(w, r, "/", http.StatusSeeOther)
			return nil
		}

		data := LoginTemplateData{"The password you entered was invalid!"}
		err = document.Execute(w, data)
		if err != nil {
			return HTTPerror{500}
		}
		return nil
	}

	err = document.Execute(w, nil)
	if err != nil {
		return HTTPerror{500}
	}

	return nil
}

func (p *Pages) signup(w http.ResponseWriter, r *http.Request, user_details *UserDetails) ErrorResponse {
	fmt.Println("Called signup")

	if r.Method == "POST" {
		stmt, err := p.db.Prepare("INSERT INTO user_data (Username, Password, Email, DOB, FirstName, LastName) VALUES (?, ?, ?, ?, ?, ?)")
		if err != nil {
			return HTTPerror{500}
		}

		err = r.ParseForm()
		if err != nil {
			return HTTPerror{500}
		}

		DOB := r.PostForm["dob-year"][0] + "-" + r.PostForm["dob-month"][0] + "-" + r.PostForm["dob-day"][0]
		password_hash, err := bcrypt.GenerateFromPassword([]byte(r.PostForm["password"][0]), 12)
		if err != nil {
			return HTTPerror{500}
		}

		_, err = stmt.Exec(
			r.PostForm["username"][0],
			string(password_hash),
			r.PostForm["email"][0],
			DOB,
			r.PostForm["firstname"][0],
			r.PostForm["lastname"][0],
		)
		if err != nil {
			return HTTPerror{500}
		}

		defer stmt.Close()

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return nil
	}

	document, err := template.ParseFiles(p.template_path+"base.html", p.template_path+"signup.html")
	if err != nil {
		return HTTPerror{500}
	}

	err = document.Execute(w, nil)
	if err != nil {
		return HTTPerror{500}
	}

	return nil
}

func main() {
	pages := new(Pages)
	pages.db, _ = sql.Open("mysql", "matthew:MysqlPassword111@tcp(127.0.0.1:3306)/UKIW")
	pages.template_path = "templates/"

	//testng only
	fs := http.FileServer(http.Dir("/home/matthew/go/src/UKIWcoursework/static"))
	http.Handle("/static/", http.StripPrefix("/static", fs))

	http.Handle("/", handler(pages.home))
	http.Handle("/signup", handler(pages.signup))
	http.Handle("/login", handler(pages.login))

	http.ListenAndServe("192.168.1.105:8000", nil)
}
