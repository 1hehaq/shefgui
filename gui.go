package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/diamondburned/gotk4/pkg/gio/v2"
	"github.com/diamondburned/gotk4/pkg/glib/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

type App struct {
	app             *gtk.Application
	window          *gtk.ApplicationWindow
	builder         *gtk.Builder
	queryEntry      *gtk.Entry
	facetCombo      *gtk.ComboBoxText
	resultsList     *gtk.ListBox
	resultsCount    *gtk.Label
	spinner         *gtk.Spinner
	searchBtn       *gtk.Button
	copyBtn         *gtk.Button
	downloadBtn     *gtk.Button
	lastResults     []string
	originalResults []string
	preferences     *Preferences
}

type Preferences struct {
	DefaultFacet string
}

func NewApp() *App {
	return &App{
		app: gtk.NewApplication("com.shef.app", gio.ApplicationFlagsNone),
	}
}

func (a *App) Run() int {
	a.app.ConnectActivate(a.activate)
	return a.app.Run(nil)
}

func (a *App) activate() {
	a.builder = gtk.NewBuilderFromFile("window.ui")

	a.window = a.builder.GetObject("window").Cast().(*gtk.ApplicationWindow)
	a.window.SetApplication(a.app)

	a.queryEntry = a.builder.GetObject("query_entry").Cast().(*gtk.Entry)
	a.facetCombo = a.builder.GetObject("facet_combo").Cast().(*gtk.ComboBoxText)
	a.resultsList = a.builder.GetObject("results_list").Cast().(*gtk.ListBox)
	a.resultsCount = a.builder.GetObject("results_count").Cast().(*gtk.Label)
	a.spinner = a.builder.GetObject("spinner").Cast().(*gtk.Spinner)
	a.searchBtn = a.builder.GetObject("search_button").Cast().(*gtk.Button)
	a.copyBtn = a.builder.GetObject("copy_button").Cast().(*gtk.Button)
	a.downloadBtn = a.builder.GetObject("download_button").Cast().(*gtk.Button)

	a.setupFacets()
	a.connectSignals()
	a.setupMenuActions()
	a.loadPreferences()
	a.resultsCount.SetVisible(false)
	a.window.Present()
}

func (a *App) setupFacets() {
	for _, facet := range shodanFacets {
		a.facetCombo.AppendText(facet)
	}
	a.facetCombo.SetActive(0)
}

func (a *App) connectSignals() {
	a.searchBtn.ConnectClicked(func() {
		a.performSearch()
	})

	titleSearchButton := a.builder.GetObject("title_search_button").Cast().(*gtk.Button)
	titleSearchButton.ConnectClicked(func() {
		a.showSimpleSearch()
	})

	a.queryEntry.ConnectActivate(func() {
		a.performSearch()
	})

	a.copyBtn.ConnectClicked(func() {
		a.copyResults()
	})

	a.downloadBtn.ConnectClicked(func() {
		a.downloadResults()
	})
}

func (a *App) performSearch() {
	a.onSearch()
}

func (a *App) onSearch() {
	query := a.queryEntry.Text()
	if strings.TrimSpace(query) == "" {
		return
	}

	facet := a.facetCombo.ActiveText()

	a.searchBtn.SetSensitive(false)
	a.spinner.SetVisible(true)
	a.spinner.Start()

	go func() {
		results, err := searchShodan(query, facet)

		glib.IdleAdd(func() {
			a.spinner.Stop()
			a.spinner.SetVisible(false)
			a.searchBtn.SetSensitive(true)

			if err != nil {
				a.showError(err.Error())
				return
			}

			a.displayResults(results)
		})
	}()
}

func (a *App) displayResults(results []string) {
	a.clearResults()
	a.lastResults = results
	a.resultsCount.SetText(fmt.Sprintf("%d", len(results)))
	a.resultsCount.SetVisible(true)

	for _, result := range results {
		row := gtk.NewListBoxRow()
		label := gtk.NewLabel(result)
		label.SetHAlign(gtk.AlignStart)
		label.SetMarginTop(8)
		label.SetMarginBottom(8)
		label.SetMarginStart(12)
		label.SetMarginEnd(12)
		label.SetSelectable(true)
		row.SetChild(label)
		a.resultsList.Append(row)
	}
}

func (a *App) showError(message string) {
	a.clearResults()
	a.lastResults = nil
	a.resultsCount.SetVisible(false)

	row := gtk.NewListBoxRow()
	label := gtk.NewLabel("Error: " + message)
	label.SetHAlign(gtk.AlignStart)
	label.SetMarginTop(8)
	label.SetMarginBottom(8)
	label.SetMarginStart(12)
	label.SetMarginEnd(12)
	label.AddCSSClass("error")
	row.SetChild(label)
	a.resultsList.Append(row)
}

func (a *App) clearResults() {
	for {
		row := a.resultsList.RowAtIndex(0)
		if row == nil {
			break
		}
		a.resultsList.Remove(row)
	}
}

func (a *App) copyResults() {
	if len(a.lastResults) == 0 {
		return
	}

	clipboard := a.window.Clipboard()
	text := strings.Join(a.lastResults, "\n")
	clipboard.SetText(text)
}

func (a *App) downloadResults() {
	if len(a.lastResults) == 0 {
		return
	}

	dialog := gtk.NewFileDialog()
	dialog.SetTitle("Save shef Results")
	dialog.SetInitialName("shef.txt")

	homeDir, _ := os.UserHomeDir()
	if homeDir != "" {
		file := gio.NewFileForPath(homeDir)
		dialog.SetInitialFolder(file)
	}

	dialog.Save(context.Background(), &a.window.Window, func(result gio.AsyncResulter) {
		file, err := dialog.SaveFinish(result)
		if err != nil {
			return
		}
		if file != nil {
			path := file.Path()
			text := strings.Join(a.lastResults, "\n")
			writeErr := os.WriteFile(path, []byte(text), 0644)
			if writeErr != nil {
				fmt.Printf("Error saving file: %v\n", writeErr)
			} else {
				fmt.Printf("Results saved to: %s\n", path)
			}
		}
	})
}

func (a *App) setupMenuActions() {
	preferencesAction := gio.NewSimpleAction("preferences", nil)
	preferencesAction.ConnectActivate(func(parameter *glib.Variant) {
		a.showPreferences()
	})
	a.app.AddAction(preferencesAction)

	aboutAction := gio.NewSimpleAction("about", nil)
	aboutAction.ConnectActivate(func(parameter *glib.Variant) {
		a.showAbout()
	})
	a.app.AddAction(aboutAction)

	helpAction := gio.NewSimpleAction("help", nil)
	helpAction.ConnectActivate(func(parameter *glib.Variant) {
		a.showHelp()
	})
	a.app.AddAction(helpAction)
}

func (a *App) showPreferences() {
	dialog := gtk.NewDialog()
	dialog.SetTitle("Preferences")
	dialog.SetModal(true)
	dialog.SetTransientFor(&a.window.Window)
	dialog.SetDefaultSize(420, 200)

	contentArea := dialog.ContentArea()
	contentArea.SetSpacing(0)
	contentArea.SetMarginTop(24)
	contentArea.SetMarginBottom(24)
	contentArea.SetMarginStart(24)
	contentArea.SetMarginEnd(24)

	mainBox := gtk.NewBox(gtk.OrientationVertical, 20)
	mainBox.SetHExpand(true)

	headerLabel := gtk.NewLabel("Search Settings")
	headerLabel.AddCSSClass("heading")
	headerLabel.SetHAlign(gtk.AlignStart)
	headerLabel.SetMarginBottom(8)

	settingsFrame := gtk.NewFrame("")
	settingsFrame.AddCSSClass("card")
	settingsFrame.SetMarginTop(0)
	settingsFrame.SetMarginBottom(0)

	frameBox := gtk.NewBox(gtk.OrientationVertical, 16)
	frameBox.SetMarginTop(20)
	frameBox.SetMarginBottom(20)
	frameBox.SetMarginStart(20)
	frameBox.SetMarginEnd(20)

	facetBox := gtk.NewBox(gtk.OrientationHorizontal, 16)
	facetBox.SetHExpand(true)

	facetLabelBox := gtk.NewBox(gtk.OrientationVertical, 4)
	facetLabel := gtk.NewLabel("Default Search Facet")
	facetLabel.SetHAlign(gtk.AlignStart)
	facetLabel.AddCSSClass("body")
	
	facetDescription := gtk.NewLabel("Choose the default facet for new searches")
	facetDescription.SetHAlign(gtk.AlignStart)
	facetDescription.AddCSSClass("caption")
	facetDescription.AddCSSClass("dim-label")
	
	facetLabelBox.Append(facetLabel)
	facetLabelBox.Append(facetDescription)

	facetCombo := gtk.NewComboBoxText()
	facetCombo.SetHAlign(gtk.AlignEnd)
	facetCombo.SetVAlign(gtk.AlignCenter)
	facetCombo.SetSizeRequest(180, -1)
	
	facets := getFacets()
	for _, facet := range facets {
		facetCombo.AppendText(facet)
	}
	
	for i, facet := range facets {
		if facet == a.preferences.DefaultFacet {
			facetCombo.SetActive(i)
			break
		}
	}
	if facetCombo.Active() == -1 {
		facetCombo.SetActive(0)
	}

	facetBox.Append(facetLabelBox)
	facetBox.Append(facetCombo)

	frameBox.Append(facetBox)
	settingsFrame.SetChild(frameBox)

	mainBox.Append(headerLabel)
	mainBox.Append(settingsFrame)

	contentArea.Append(mainBox)

	dialog.AddButton("Cancel", int(gtk.ResponseCancel))
	dialog.AddButton("Save", int(gtk.ResponseAccept))

	dialog.ConnectResponse(func(responseId int) {
		if responseId == int(gtk.ResponseAccept) {
			a.preferences.DefaultFacet = facetCombo.ActiveText()
			facets := getFacets()
			for i, facet := range facets {
				if facet == a.preferences.DefaultFacet {
					a.facetCombo.SetActive(i)
					break
				}
			}
			fmt.Println("Default facet set to:", a.preferences.DefaultFacet)
		}
		dialog.Close()
	})

	dialog.Show()
}

func (a *App) loadPreferences() {
	a.preferences = &Preferences{
		DefaultFacet: "product",
	}

	facets := getFacets()
	for i, facet := range facets {
		if facet == a.preferences.DefaultFacet {
			a.facetCombo.SetActive(i)
			break
		}
	}
}

func (a *App) savePreferences() {
	fmt.Printf("Default facet: %s\n", a.preferences.DefaultFacet)
}

func (a *App) showAbout() {
	builder := gtk.NewBuilderFromFile("about.ui")
	aboutDialog := builder.GetObject("about_dialog").Cast().(*gtk.AboutDialog)
	aboutDialog.SetTransientFor(&a.window.Window)
	aboutDialog.Show()
}

func (a *App) showFind() {
	if len(a.lastResults) == 0 {
		return
	}

	builder := gtk.NewBuilderFromFile("find.ui")
	dialog := builder.GetObject("find_dialog").Cast().(*gtk.Dialog)
	dialog.SetTransientFor(&a.window.Window)

	findEntry := builder.GetObject("find_entry").Cast().(*gtk.Entry)
	caseSensitiveCheck := builder.GetObject("case_sensitive_check").Cast().(*gtk.CheckButton)
	matchCountLabel := builder.GetObject("match_count_label").Cast().(*gtk.Label)
	findPrevButton := builder.GetObject("find_prev_button").Cast().(*gtk.Button)
	findNextButton := builder.GetObject("find_next_button").Cast().(*gtk.Button)
	findCloseButton := builder.GetObject("find_close_button").Cast().(*gtk.Button)

	currentMatch := -1
	matches := []int{}

	updateMatches := func() {
		query := findEntry.Text()
		if query == "" {
			matches = []int{}
			matchCountLabel.SetText("")
			return
		}

		matches = []int{}
		caseSensitive := caseSensitiveCheck.Active()

		for i, result := range a.lastResults {
			searchText := result
			searchQuery := query

			if !caseSensitive {
				searchText = strings.ToLower(searchText)
				searchQuery = strings.ToLower(searchQuery)
			}

			if strings.Contains(searchText, searchQuery) {
				matches = append(matches, i)
			}
		}

		if len(matches) > 0 {
			matchCountLabel.SetText(fmt.Sprintf("%d matches", len(matches)))
		} else {
			matchCountLabel.SetText("No matches")
		}
		currentMatch = -1
	}

	findEntry.ConnectChanged(func() {
		updateMatches()
	})

	caseSensitiveCheck.ConnectToggled(func() {
		updateMatches()
	})

	findNextButton.ConnectClicked(func() {
		if len(matches) > 0 {
			currentMatch = (currentMatch + 1) % len(matches)
			fmt.Printf("Next match: result %d\n", matches[currentMatch])
		}
	})

	findPrevButton.ConnectClicked(func() {
		if len(matches) > 0 {
			if currentMatch <= 0 {
				currentMatch = len(matches) - 1
			} else {
				currentMatch--
			}
			fmt.Printf("Previous match: result %d\n", matches[currentMatch])
		}
	})

	findCloseButton.ConnectClicked(func() {
		dialog.Close()
	})

	findEntry.ConnectActivate(func() {
		if len(matches) > 0 {
			currentMatch = (currentMatch + 1) % len(matches)
			fmt.Printf("Enter pressed: result %d\n", matches[currentMatch])
		}
	})

	dialog.Show()
	findEntry.GrabFocus()
}

func (a *App) showHelp() {
	builder := gtk.NewBuilderFromFile("help.ui")
	dialog := builder.GetObject("help_dialog").Cast().(*gtk.Dialog)
	dialog.SetTransientFor(&a.window.Window)

	facetsFlowBox := builder.GetObject("facets_flowbox").Cast().(*gtk.FlowBox)
	searchEntry := builder.GetObject("facet_search_entry").Cast().(*gtk.SearchEntry)

	allFacets := a.loadFacetsFromJSON(facetsFlowBox)

	searchEntry.ConnectSearchChanged(func() {
		query := strings.ToLower(searchEntry.Text())
		a.filterFacets(facetsFlowBox, allFacets, query)
	})

	closeButton := builder.GetObject("help_close_button").Cast().(*gtk.Button)
	closeButton.ConnectClicked(func() {
		dialog.Close()
	})

	dialog.Show()
}

func (a *App) showSimpleSearch() {
	if len(a.lastResults) == 0 {
		fmt.Println("No results to search")
		return
	}

	dialog := gtk.NewDialog()
	dialog.SetTitle("Search Results")
	dialog.SetModal(true)
	dialog.SetTransientFor(&a.window.Window)
	dialog.SetDefaultSize(500, 400)

	contentArea := dialog.ContentArea()
	contentArea.SetSpacing(12)
	contentArea.SetMarginTop(16)
	contentArea.SetMarginBottom(16)
	contentArea.SetMarginStart(16)
	contentArea.SetMarginEnd(16)

	searchEntry := gtk.NewSearchEntry()
	searchEntry.SetPlaceholderText("Search in results...")
	searchEntry.SetHExpand(true)

	scrolled := gtk.NewScrolledWindow()
	scrolled.SetPolicy(gtk.PolicyAutomatic, gtk.PolicyAutomatic)
	scrolled.SetVExpand(true)
	scrolled.SetHExpand(true)

	resultsList := gtk.NewListBox()
	scrolled.SetChild(resultsList)

	contentArea.Append(searchEntry)
	contentArea.Append(scrolled)

	if a.originalResults == nil {
		a.originalResults = make([]string, len(a.lastResults))
		copy(a.originalResults, a.lastResults)
	}

	for _, result := range a.originalResults {
		row := gtk.NewListBoxRow()
		label := gtk.NewLabel(result)
		label.SetHAlign(gtk.AlignStart)
		label.SetWrap(true)
		label.SetXAlign(0)
		label.SetMarginTop(8)
		label.SetMarginBottom(8)
		label.SetMarginStart(12)
		label.SetMarginEnd(12)
		row.SetChild(label)
		resultsList.Append(row)
	}

	searchEntry.ConnectSearchChanged(func() {
		query := strings.TrimSpace(strings.ToLower(searchEntry.Text()))
		
		resultsList.RemoveAll()
		
		for _, result := range a.originalResults {
			if query == "" || strings.Contains(strings.ToLower(result), query) {
				row := gtk.NewListBoxRow()
				label := gtk.NewLabel(result)
				label.SetHAlign(gtk.AlignStart)
				label.SetWrap(true)
				label.SetXAlign(0)
				label.SetMarginTop(8)
				label.SetMarginBottom(8)
				label.SetMarginStart(12)
				label.SetMarginEnd(12)
				row.SetChild(label)
				resultsList.Append(row)
			}
		}
	})

	dialog.AddButton("Close", int(gtk.ResponseClose))

	dialog.ConnectResponse(func(responseId int) {
		dialog.Close()
	})

	dialog.Show()
	searchEntry.GrabFocus()
}

type FacetData struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type FacetsJSON struct {
	Facets []FacetData `json:"facets"`
}

func (a *App) loadFacetsFromJSON(flowBox *gtk.FlowBox) []FacetData {
	data, err := os.ReadFile("facets.json")
	if err != nil {
		fmt.Printf("Error reading facets.json: %v\n", err)
		return nil
	}

	var facetsData FacetsJSON
	err = json.Unmarshal(data, &facetsData)
	if err != nil {
		fmt.Printf("Error parsing facets.json: %v\n", err)
		return nil
	}

	for _, facet := range facetsData.Facets {
		child := gtk.NewFlowBoxChild()

		facetBox := gtk.NewBox(gtk.OrientationVertical, 6)
		facetBox.SetMarginTop(12)
		facetBox.SetMarginBottom(12)
		facetBox.SetMarginStart(12)
		facetBox.SetMarginEnd(12)
		facetBox.SetSizeRequest(200, -1)

		nameLabel := gtk.NewLabel(facet.Name)
		nameLabel.SetHAlign(gtk.AlignStart)
		nameLabel.AddCSSClass("heading")

		descLabel := gtk.NewLabel(facet.Description)
		descLabel.SetHAlign(gtk.AlignStart)
		descLabel.SetWrap(true)
		descLabel.SetXAlign(0)
		descLabel.SetMaxWidthChars(30)
		descLabel.AddCSSClass("caption")

		facetBox.Append(nameLabel)
		facetBox.Append(descLabel)

		child.SetChild(facetBox)
		flowBox.Append(child)
	}

	return facetsData.Facets
}

func (a *App) filterResults(query string) {
	a.clearResults()
	
	query = strings.TrimSpace(strings.ToLower(query))
	
	if query == "" {
		a.restoreOriginalResults()
		return
	}
	
	var filteredResults []string
	for _, result := range a.originalResults {
		if strings.Contains(strings.ToLower(result), query) {
			filteredResults = append(filteredResults, result)
		}
	}
	
	a.lastResults = filteredResults
	a.resultsCount.SetText(fmt.Sprintf("%d", len(filteredResults)))
	
	for _, result := range filteredResults {
		row := gtk.NewListBoxRow()
		label := gtk.NewLabel(result)
		label.SetHAlign(gtk.AlignStart)
		label.SetWrap(true)
		label.SetXAlign(0)
		label.SetMarginTop(8)
		label.SetMarginBottom(8)
		label.SetMarginStart(12)
		label.SetMarginEnd(12)
		row.SetChild(label)
		a.resultsList.Append(row)
	}
}

func (a *App) restoreOriginalResults() {
	if a.originalResults == nil {
		return
	}
	
	a.lastResults = make([]string, len(a.originalResults))
	copy(a.lastResults, a.originalResults)
	
	a.clearResults()
	a.resultsCount.SetText(fmt.Sprintf("%d", len(a.lastResults)))
	
	for _, result := range a.lastResults {
		row := gtk.NewListBoxRow()
		label := gtk.NewLabel(result)
		label.SetHAlign(gtk.AlignStart)
		label.SetWrap(true)
		label.SetXAlign(0)
		label.SetMarginTop(8)
		label.SetMarginBottom(8)
		label.SetMarginStart(12)
		label.SetMarginEnd(12)
		row.SetChild(label)
		a.resultsList.Append(row)
	}
	
	a.originalResults = nil
}

func (a *App) filterFacets(flowBox *gtk.FlowBox, allFacets []FacetData, query string) {
	flowBox.RemoveAll()

	for _, facet := range allFacets {
		if query == "" ||
			strings.Contains(strings.ToLower(facet.Name), query) ||
			strings.Contains(strings.ToLower(facet.Description), query) {

			child := gtk.NewFlowBoxChild()

			facetBox := gtk.NewBox(gtk.OrientationVertical, 6)
			facetBox.SetMarginTop(12)
			facetBox.SetMarginBottom(12)
			facetBox.SetMarginStart(12)
			facetBox.SetMarginEnd(12)
			facetBox.SetSizeRequest(200, -1)

			nameLabel := gtk.NewLabel(facet.Name)
			nameLabel.SetHAlign(gtk.AlignStart)
			nameLabel.AddCSSClass("heading")

			descLabel := gtk.NewLabel(facet.Description)
			descLabel.SetHAlign(gtk.AlignStart)
			descLabel.SetWrap(true)
			descLabel.SetXAlign(0)
			descLabel.SetMaxWidthChars(30)
			descLabel.AddCSSClass("caption")

			facetBox.Append(nameLabel)
			facetBox.Append(descLabel)

			child.SetChild(facetBox)
			flowBox.Append(child)
		}
	}
}

func getFacets() []string {
	return []string{
		"asn", "bitcoin.ip", "bitcoin.ip_count", "bitcoin.port", "bitcoin.user_agent",
		"bitcoin.version", "city", "cloud.provider", "cloud.region", "cloud.service",
		"country", "cpe", "device", "domain", "has_screenshot", "hash",
		"http.component", "http.component_category", "http.dom_hash", "http.favicon.hash",
		"http.headers_hash", "http.html_hash", "http.robots_hash", "http.server_hash",
		"http.status", "http.title", "http.title_hash", "http.waf", "ip", "isp",
		"link", "mongodb.database.name", "ntp.ip", "ntp.ip_count", "ntp.more",
		"ntp.port", "org", "os", "port", "postal", "product", "redis.key",
		"region", "rsync.module", "screenshot.hash", "screenshot.label",
		"snmp.contact", "snmp.location", "snmp.name", "ssh.cipher", "ssh.fingerprint",
		"ssh.hassh", "ssh.mac", "ssh.type", "ssl.alpn", "ssl.cert.alg",
		"ssl.cert.expired", "ssl.cert.extension", "ssl.cert.fingerprint",
		"ssl.cert.issuer.cn", "ssl.cert.pubkey.bits", "ssl.cert.pubkey.type",
		"ssl.cert.serial", "ssl.cert.subject.cn", "ssl.chain_count",
		"ssl.cipher.bits", "ssl.cipher.name", "ssl.cipher.version", "ssl.ja3s",
		"ssl.jarm", "ssl.version", "state", "tag", "telnet.do", "telnet.dont",
		"telnet.option", "telnet.will", "telnet.wont", "uptime", "version",
		"vuln", "vuln.verified",
	}
}
