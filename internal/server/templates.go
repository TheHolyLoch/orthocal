// Orthocal - Developed by dgm (dgm@tuta.com)
// orthocal/internal/server/templates.go

package server

const templates = `
{{define "header"}}
<!doctype html>
<html lang="en">
<head>
	<meta charset="utf-8">
	<meta name="viewport" content="width=device-width, initial-scale=1">
	<title>{{.Title}} - Orthocal</title>
	{{if .StyleSheet}}<style>{{.StyleSheet}}</style>{{else}}<link rel="stylesheet" href="{{.AssetPrefix}}style.css">{{end}}
</head>
<body>
	<main class="page">
		<header class="topbar">
			<a class="brand" href="{{.HomeLink}}">Orthocal</a>
			<nav class="nav" aria-label="Primary">
				{{if .DateValue}}<a href="{{.DayPrefix}}{{.DateValue}}{{.LinkSuffix}}">Day</a>{{end}}
				<a href="{{.DayPrefix}}{{.Today}}{{.LinkSuffix}}">Today</a>
				{{if .DateValue}}<a href="{{.SaintsPrefix}}{{.DateValue}}{{.LinkSuffix}}">Saints</a>{{end}}
				{{if .DateValue}}<a href="{{.ReadPrefix}}{{.DateValue}}{{.LinkSuffix}}">Readings</a>{{end}}
				{{if .DateValue}}<a href="{{.HymnsPrefix}}{{.DateValue}}{{.LinkSuffix}}">Hymns</a>{{end}}
				{{if .DateValue}}<a href="{{.APIPrefix}}{{.DateValue}}{{.APISuffix}}">API JSON</a>{{end}}
			</nav>
		</header>
{{end}}

{{define "footer"}}
	</main>
	{{if .AppScript}}<script>{{.AppScript}}</script>{{else}}<script src="{{.AssetPrefix}}app.js"></script>{{end}}
</body>
</html>
{{end}}

{{define "day_page"}}
{{template "header" .}}
<section class="hero">
	<div class="knot"></div>
	<h1>{{.DayView.Day.DataHeader}}</h1>
	{{if .DayView.Day.HeaderHeader}}<p class="subhead">{{.DayView.Day.HeaderHeader}}</p>{{end}}
	{{if .DayView.Day.FastingRule}}<p class="fasting">{{.DayView.Day.FastingRule}}</p>{{end}}
	<div class="controls">
		<a class="button" href="{{.DayPrefix}}{{.PrevDate}}{{.LinkSuffix}}">Previous Day</a>
		<a class="button" href="{{.DayPrefix}}{{.NextDate}}{{.LinkSuffix}}">Next Day</a>
		<a class="button" href="{{.DayPrefix}}{{.Today}}{{.LinkSuffix}}">Today</a>
		<label>Date
			<input data-date-picker type="date" value="{{.DateValue}}">
		</label>
	</div>
</section>

<section class="grid">
	{{if .DayView.Saints}}
	<article class="panel">
		<h2>Saints</h2>
		<ol>
		{{range .DayView.Saints}}
			<li>
				{{if .ServiceRankCode}}[{{.ServiceRankCode}}] {{end}}{{.Name}}
				{{if .IsPrimary}}<span class="tag">primary</span>{{end}}
				{{if .IsWestern}}<span class="tag">western</span><span class="tag">Britain and Ireland</span>{{end}}
			</li>
		{{end}}
		</ol>
	</article>
	{{end}}

	{{if .DayView.ScriptureReadings}}
	<article class="panel">
		<h2>Scripture</h2>
		<ul>
		{{range .DayView.ScriptureReadings}}
			<li>{{.VerseReference}}{{if .Description}} - {{.Description}}{{end}}</li>
		{{end}}
		</ul>
	</article>
	{{end}}

	{{if .HymnCount}}
	<article class="panel">
		<h2>Hymns</h2>
		<p>{{.HymnCount}} hymns available.</p>
		<a href="{{.HymnsPrefix}}{{.DateValue}}{{.LinkSuffix}}">View hymns</a>
	</article>
	{{end}}
</section>
{{template "footer" .}}
{{end}}

{{define "saints_page"}}
{{template "header" .}}
<section class="hero">
	<div class="knot"></div>
	<h1>Saints</h1>
	<p class="subhead">{{.SaintsView.Day.DataHeader}}</p>
</section>
<section class="panel">
	{{if .SaintsView.Saints}}
	<ol>
	{{range .SaintsView.Saints}}
		<li>
			{{if .ServiceRankCode}}[{{.ServiceRankCode}}] {{end}}{{.Name}}
			{{if .IsPrimary}}<span class="tag">primary</span>{{end}}
			{{if .IsWestern}}<span class="tag">western</span>{{end}}
		</li>
	{{end}}
	</ol>
	{{else}}
	<p>No saints found.</p>
	{{end}}
</section>
{{template "footer" .}}
{{end}}

{{define "readings_page"}}
{{template "header" .}}
<section class="hero">
	<div class="knot"></div>
	<h1>Scripture Readings</h1>
	<p class="subhead">{{.ReadingsView.Day.DataHeader}}</p>
</section>
<section class="panel">
	{{if .ReadingsView.ScriptureReadings}}
	<ol>
	{{range .ReadingsView.ScriptureReadings}}
		<li>{{.VerseReference}}{{if .Description}} - {{.Description}}{{end}}</li>
	{{end}}
	</ol>
	{{else}}
	<p>No scripture readings found.</p>
	{{end}}
</section>
{{template "footer" .}}
{{end}}

{{define "hymns_page"}}
{{template "header" .}}
<section class="hero">
	<div class="knot"></div>
	<h1>Hymns</h1>
	<p class="subhead">{{.HymnsView.Day.DataHeader}}</p>
</section>
<section class="panel">
	{{if .HymnsView.Hymns}}
	{{range .HymnsView.Hymns}}
	<article class="hymn">
		<h2>{{.Title}}</h2>
		<div class="hymn-meta">{{.HymnType}}{{if .Tone}} - Tone {{.Tone}}{{end}}</div>
		{{if .Text}}<p>{{.Text}}</p>{{end}}
	</article>
	{{end}}
	{{else}}
	<p>No hymns found.</p>
	{{end}}
</section>
{{template "footer" .}}
{{end}}

{{define "not_found_page"}}
{{template "header" .}}
<section class="hero">
	<div class="knot"></div>
	<h1>Date Not Found</h1>
	<p class="subhead">{{.Error}}</p>
	<div class="controls">
		<a class="button" href="{{.DayPrefix}}{{.Today}}{{.LinkSuffix}}">Today</a>
		<label>Date
			<input data-date-picker type="date" value="{{.DateValue}}">
		</label>
	</div>
</section>
{{template "footer" .}}
{{end}}

{{define "error_page"}}
{{template "header" .}}
<section class="error">
	<h1>Error</h1>
	<p>{{.Error}}</p>
</section>
{{template "footer" .}}
{{end}}
`
