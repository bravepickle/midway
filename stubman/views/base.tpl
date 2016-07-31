{{ define "base" }}<!DOCTYPE html>
<html lang="en">
<head>
	<!-- bootstrap -->
	<meta name="viewport" content="width=device-width, initial-scale=1">
	<link rel="stylesheet" href="static/css/bootstrap.min.css" />
	<link rel="stylesheet" href="static/css/bootstrap-theme.min.css" />
	<script src="static/js/bootstrap.min.js" crossorigin="anonymous"></script>

    {{ template "title" . }}
</head>
<body>
    {{ template "scripts" . }}
    {{ template "sidebar" . }}
	<div class="container">
    {{ template "content" . }}
	</div>
<footer>
    Beta Version
</footer>
</body>
</html>
{{ end }}
// We define empty blocks for optional content so we don't have to define a block in child templates that don't need them
{{ define "scripts" }}{{ end }}
{{ define "sidebar" }}{{ end }}