<!DOCTYPE html>
<html>
<head>
	<meta charset="UTF-8">
	<title>Frage {{.Qno}} – Umfrage der CSG-Abizeitung 2017</title>
</head>
<body>
	<form action="/q/{{.Qno}}" method="post">
		<h3>Frage {{.Qno}}. {{.String}}</h3>

		<input type="text" name="answer"><br>

		<input type="submit" value="OK">
	</form>
</body>
</html>
