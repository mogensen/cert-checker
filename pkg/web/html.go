package web

import (
	"fmt"
	"html/template"
	"io"
	"time"

	"github.com/mogensen/cert-checker/pkg/models"
)

// Colors used https://loading.io/color/feature/evaneos/
// HTML and CSS inspired by https://github.com/bderenzo/tinystatus/

const createTpl = `
<!DOCTYPE html><html lang="en">
<head>
	<meta charset="utf-8">
	<meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">
	<meta name="color-scheme" content="dark light">
	<meta http-equiv="refresh" content="30">
	<title>cert-checker</title>
	<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/github-fork-ribbon-css/0.2.3/gh-fork-ribbon.min.css" />

	<script>
		// function to set a given theme/color-scheme
		function setTheme(themeName) {
			localStorage.setItem('theme', themeName);
			document.documentElement.className = themeName;
		}

		// function to toggle between light and dark theme
		function toggleTheme() {
			if (localStorage.getItem('theme') === 'theme-dark') {
				setTheme('theme-light');
			} else {
				setTheme('theme-dark');
			}
		}

		// Immediately invoked function to set the theme on initial load
		(function () {

			var currentTheme = localStorage.getItem('theme')

			if (currentTheme === 'theme-dark') {
				setTheme('theme-dark');
				return
			} else if (currentTheme === 'theme-light') {
				setTheme('theme-light');
				return
			}
			// Default to reading the browser preference
			if (window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches) {
				setTheme('theme-dark');
			}
		})();

		document.addEventListener("DOMContentLoaded", function() {
			document.getElementById('metrics-link').href =  window.location.protocol + '//' + window.location.hostname+':8080/metrics'; // TODO Template this
		})

	</script>


	<style>

		.theme-light {
			--color-background: #fbfbfe;
			--font-color: #000000;
			--color-secondary: #4ba6f5;
			--color-success: #33cc99;
			--color-yellow: #fdca30;
			--color-warning: #f79400;
			--color-failed: #f34235;
		}
		.theme-dark {
			--color-background: #243133;
			--font-color: #ffffff;
			--color-secondary: #4ba6f5;
			--color-success: #33cc99;
			--color-yellow: #fdca30;
			--color-warning: #f79400;
			--color-failed: #f34235;
		}
		.container button {
			float: right;
			color: var(--font-color);
			background: #4ba6f5;
			padding: 10px 20px;
			border: 0;
			border-radius: 5px;
		  }

		  ul.footer {
			display: inline-grid;
			grid-auto-flow: row;
			grid-gap: 24px;
			justify-items: center;
			margin: auto;
		  }

		  @media (min-width: 500px) {
			ul.footer {
			  grid-auto-flow: column;
			}
		  }

		  a {
			color: #7b7b7b;
			text-decoration: none;
		  }

		  div.footer {
			display: flex;
			width: 100%;
			line-height: 1.3;
		  }

		body { font-family: segoe ui,Roboto,Oxygen-Sans,Ubuntu,Cantarell,helvetica neue,Verdana,sans-serif; background: var(--color-background);  color: var(--font-color);}
		h1 { margin-top: 30px; }
		h2 { margin-top: 30px; }
		h3 { margin-top: 5px;  margin-bottom: 5px; }
		ul { padding: 0px; }
		li { list-style: none; margin-bottom: 2px; padding: 5px; }
		.lastupdate{ color: var(--color-success); }
		.container { max-width: 900px; width: 100%; margin: 15px auto; }
		.panel { text-align: center; padding: 10px; border: 0px; border-radius: 5px; }
		.failed-bg  { color: white; background-color: var(--color-failed); }
		.warning-bg  { color: white; background-color: var(--color-warning); }
		.success-bg { color: white; background-color: var(--color-success) }
		.hidden  { display: none; }
		.failed  { color: var(--color-failed); }
		.warning  { color: var(--color-warning); }
		.success { color: var(--color-success) }
		.small { font-size: 80%; }
		.status { float: right; }
		.github-fork-ribbon:before {background-color: #333;}
	</style>
</head>
<body>
<div class="container">
	<button id="switch" onclick="toggleTheme()">Dark mode</button>
</div>
<a class="github-fork-ribbon" href="https://github.com/mogensen/cert-checker" data-ribbon="Fork me on GitHub" title="Fork me on GitHub">Fork me on GitHub</a>
<div class='container'>
<h1>cert-checker</h1>

{{ if .BadCerts }}
	<ul><li class='panel failed-bg'>{{ len .BadCerts }} Bad certificate(s)</li></ul>
{{ else }}
	<ul><li class='panel success-bg'>All Systems Operational</li></ul>
{{ end }}

<span class="lastupdate">Last update: {{ timeNow }}</span>

<h2>Certificates</h2>

<ul>
{{ range $key, $cert := .BadCerts }}
	<hr />
	<h3>{{ $cert.DNS }}</h3>

	<li>
	<span class='small'>
	{{ if $cert.Info.Issuer }}
		{{ $cert.Info.Issuer }}
	{{else}}
		No valid issuer fould
	{{end}}
	</span>

	<span class='small {{ expireTimeToCSS $cert.Info.NotAfter }}'>Expires: {{ ts $cert.Info.NotAfter }},</span>

	<span class='small {{ tlsToCSS $cert.Info.MinimumTLSVersion }}'>{{ $cert.Info.MinimumTLSVersion }}</span>

	{{ if $cert.Info.Error }}
		<span class='status failed'>Error</span>
		<li class='panel failed-bg'>{{ $cert.Info.Error }}</li>
	{{ else }}
		<span class='status success'>All Good!</span>
	{{ end }}

	</li>

{{- end}}

{{ range $key, $cert := .GoodCerts }}
	<hr />
	<h3>{{ $cert.DNS }}</h3>

	<li>
	{{ if $cert.Info.Issuer }}
		{{ $cert.Info.Issuer }}
	{{else}}
		No valid issuer fould
	{{end}}

	<span class='small {{ expireTimeToCSS $cert.Info.NotAfter }}'>Expires: {{ ts $cert.Info.NotAfter }},</span>

	<span class='small {{ tlsToCSS $cert.Info.MinimumTLSVersion }}'>{{ $cert.Info.MinimumTLSVersion }}</span>

	{{ if $cert.Info.Error }}
		<span class='status failed'>Error</span>
		<li class='panel failed-bg'>{{ $cert.Info.Error }}</li>
	{{else}}
		<span class='status success'>All Good!</span>
	{{end}}


	</li>

{{- end}}

</ul>
</div>


<div class="footer">
  <ul class="footer">
    <li><a id="metrics-link">Metrics</a></li>
    <li><a href="https://github.com/mogensen/cert-checker/issues">Issues</a></li>
    <li><a href="https://github.com/mogensen/cert-checker">Github</a></li>
  </ul>
</div>

</body></html>
`

// templateHTML generates an html representation for the given certs, and writes the result to the io.writer
func templateHTML(certs []models.Certificate, w io.Writer) error {

	sum := internalSummery{}

	for _, c := range certs {
		if c.Info == nil {
			// Cert is not proccessed yet
			continue
		} else if c.Info.Error != "" {
			sum.BadCerts = append(sum.BadCerts, c)
		} else {
			sum.GoodCerts = append(sum.GoodCerts, c)
		}
	}

	t := template.Must(template.New("create").Funcs(
		template.FuncMap{
			"ts": func(ts string) string {
				parsedTime := pTime(ts)
				if parsedTime == nil {
					return ""
				}
				return parsedTime.Format("2006-01-02")
			},
			"timeNow": func() string {
				return time.Now().Format("2006-01-02 15:04:05")
			},
			"expiresSoon": func(ts string) bool {
				parsedTime := pTime(ts)
				if parsedTime == nil {
					return false
				}
				return parsedTime.Before(time.Now().AddDate(0, 0, 25))
			},
			"tlsToCSS": func(ts string) string {
				switch ts {
				case "SSLv3":
				case "TLS 1.0":
					return "failed"
				case "TLS 1.1":
					return "warning"
				case "TLS 1.2":
					return "success"
				case "TLS 1.3":
					return "success"
				}
				return "failed"
			},
			"expireTimeToCSS": func(ts string) string {
				parsedTime := pTime(ts)
				if parsedTime == nil {
					return "hidden"
				}
				if parsedTime.Before(time.Now()) {
					return "failed"
				}
				if parsedTime.Before(time.Now().AddDate(0, 0, 25)) {
					return "warning"
				}
				return "success"
			},
		},
	).Parse(createTpl))
	err := t.Execute(w, sum)
	if err != nil {
		return err
	}
	return nil
}

func pTime(s string) *time.Time {
	parsedTime, err := time.Parse("2006-01-02 15:04:05 -0700 MST", s)
	if err != nil || (parsedTime == time.Time{}) {
		fmt.Println(err)
		return nil
	}
	return &parsedTime
}

type internalSummery struct {
	GoodCerts    []models.Certificate
	WarningCerts []models.Certificate
	BadCerts     []models.Certificate
}
