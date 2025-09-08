package main

import (
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/diamondburned/gotk4/pkg/gdk/v4"
)

func applyStyles() {
	provider := gtk.NewCSSProvider()
	css := `
window {
	background-color: #2e3440;
	color: #d8dee9;
}

entry {
	background-color: #3b4252;
	color: #eceff4;
	border: 1px solid #4c566a;
	border-radius: 6px;
	padding: 8px;
	min-width: 200px;
}

entry:focus {
	border-color: #5e81ac;
	box-shadow: 0 0 0 2px rgba(94, 129, 172, 0.3);
}

button {
	background-color: #4c566a;
	color: #eceff4;
	border: 1px solid #5e81ac;
	border-radius: 6px;
	padding: 8px 16px;
}

button:hover {
	background-color: #5e81ac;
}

button.suggested-action {
	background-color: #5e81ac;
	color: #eceff4;
}

button.suggested-action:hover {
	background-color: #81a1c1;
}

combobox {
	background-color: #3b4252;
	color: #eceff4;
	border: 1px solid #4c566a;
	border-radius: 6px;
	min-width: 120px;
	max-width: 180px;
}

combobox button {
	background-color: #3b4252;
	border: none;
	padding: 8px;
}

combobox button:hover {
	background-color: #434c5e;
}

listbox {
	background-color: #3b4252;
	border-radius: 6px;
}

listbox row {
	background-color: #3b4252;
	color: #d8dee9;
	border-radius: 6px;
	margin: 2px;
}

listbox row:hover {
	background-color: #434c5e;
}

listbox row label {
	color: #d8dee9;
}

listbox row label.error {
	color: #bf616a;
}

.title-4 {
	color: #81a1c1;
	font-weight: bold;
}

button.flat {
	background: transparent;
	border: none;
	color: #d8dee9;
	border-radius: 6px;
	padding: 8px;
}

button.flat:hover {
	background-color: #434c5e;
	color: #eceff4;
}

button.flat:active {
	background-color: #4c566a;
}

.count-label {
	color: #81a1c1;
	font-weight: 900;
	font-size: 16px;
}

checkbutton {
	color: #d8dee9;
}

scrolledwindow {
	border: 1px solid #4c566a;
	border-radius: 6px;
}
`
	provider.LoadFromData(css)
	
	display := gdk.DisplayGetDefault()
	gtk.StyleContextAddProviderForDisplay(
		display,
		provider,
		gtk.STYLE_PROVIDER_PRIORITY_APPLICATION,
	)
}
