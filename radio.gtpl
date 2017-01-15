<!DOCTYPE html>
<html>
<head>
	<meta charset="UTF-8">
	<title>Frage {{.Qno}} – Umfrage der CSG-Abizeitung 2017</title>
	<style>
		/* Radio buttons should be ordered in four columns. */
		.radios {
			-webkit-column-count: 4; /* Chrome, Safari, Opera */
			-moz-column-count: 4;    /* Firefox */
			column-count: 4;
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

		<input type="submit" value="Weiter">
	</form>

	<form action="/q/{{.Qnoprev}}" method="get" >
		<button type="submit">Zurück</button>
	</form>
</body>
</html>
