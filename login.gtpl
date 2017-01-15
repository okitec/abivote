<!DOCTYPE html>
<html>
<head>
	<meta charset="UTF-8">
	<title>Login – Umfrage der CSG-Abizeitung 2017</title>
</head>
<body>
	<p><strong>Devnote</strong>: Kennungen und Sessions sind noch nicht implementiert. Wenn
	also zwei Leute gleichzeitig den Server benutzen, kann dieser die beiden momentan nicht
	unterscheiden. Eine Person kann also dieselbe Frage mehrfach beantworten. URL-Rewriting
	kommt noch, das wird das alles fixen.</p>

	<p>Die Statistiken werden außerdem noch nicht gespeichert.</p>

	<a href="/stats">Statistiken (nur für Auswerter, oder? Kann ich auch machen.)</a>

	<form action="/login" method="post">
		Deine Kennung: <input type="text" name="username">

		<input type="submit" value="Login">
	</form>
</body>
</html>
