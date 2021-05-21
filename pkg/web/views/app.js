// {{define "js"}}


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
    } else {
        setTheme('theme-light');
    }
})();

document.addEventListener("DOMContentLoaded", function () {
    document.getElementById('metrics-link').href = window.location.protocol + '//' + window.location.hostname + ':8080/metrics'; // TODO Template this
})

// {{ end }}
