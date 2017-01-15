<!DOCTYPE html>
<html>
<head>
	<meta charset="UTF-8">
	<title>Frage {{.Qno}} – Umfrage der CSG-Abizeitung 2017</title>
	<style>
		.radios {
			-webkit-column-count: 2; /* Chrome, Safari, Opera */
			-moz-column-count: 2;    /* Firefox */
			column-count: 2;
		}
	</style>
</head>
<body>
	<form action="/q/{{.Qno}}" method="post">
		<h3>Frage {{.Qno}}. {{.String}}</h3>

		<div class="radios">
		{{range .Choices}}
			<input type="radio" name="answer" value="{{.IdentAnswer}}"> {{.Answer}}<br>
		{{end}}
		</div>

		<input type="submit" value="OK">
	</form>
</body>
</html>
