package template

import (
	"html/template"
	"log"
	"os"
)

const Tpl = `
<!DOCTYPE html>
<html>
	<head>
		<meta charset="UTF-8">
		<title>{{.Title}}</title>
	</head>
	<body>
		<form action="/login" method="post">
			<div class="imgcontainer">
			<img src="img_avatar2.png" alt="Avatar" class="avatar">
			</div>

			<div class="container">
			<label for="uname"><b>Username</b></label>
			<input type="text" placeholder="Enter Username" name="uname" required>

			<label for="psw"><b>Password</b></label>
			<input type="password" placeholder="Enter Password" name="psw" required>

			<button type="submit">Login</button>
			<label>
				<input type="checkbox" checked="checked" name="remember"> Remember me
			</label>
			</div>

			<div class="container" style="background-color:#f1f1f1">
			<button type="button" class="cancelbtn">Cancel</button>
			<span class="psw">Forgot <a href="#">password?</a></span>
			</div>
		</form>
	</body>
</html>`

func Login() {

	check := func(err error) {
		if err != nil {
			log.Fatal(err)
		}
	}
	t, err := template.New("webpage").Parse(Tpl)
	check(err)

	data := struct {
		Title string
		Items []string
	}{
		Title: "My page",
		Items: []string{
			"My photos",
			"My blog",
		},
	}

	err = t.Execute(os.Stdout, data)
	check(err)

	noItems := struct {
		Title string
		Items []string
	}{
		Title: "My another page",
		Items: []string{},
	}

	err = t.Execute(os.Stdout, noItems)
	check(err)

}
