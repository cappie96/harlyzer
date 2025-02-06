package harlyzer

import (
	"fmt"
	"sort"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Terminal struct {
	app      *tview.Application
	table    *tview.Table
	dropdown *tview.DropDown
	input    *tview.InputField
	main     *tview.Flex
	modal    *tview.Modal
}

func NewTerminal() *Terminal {
	return new(Terminal)
}

func (t *Terminal) Init() {
	t.app = tview.NewApplication()
	t.table = tview.NewTable()
	t.dropdown = tview.NewDropDown()
	t.input = tview.NewInputField()
}

func (t *Terminal) CreateTable(har *HAR, code string, url string) {
	if har == nil || har.Log.Entries == nil {
		fmt.Println("Invalid HAR data")
		return
	}
	t.table.SetFixed(1, 1).SetBorderPadding(1, 1, 1, 1)
	t.table.SetBorder(true).SetTitle("HAR Log")

	// Headers
	headers := []string{"#", "Method", "Status", "Domain", "Url", "Server IP", "Connection", "Time (ms)"}
	t.SetTableHeader(headers)
	// Parse code filter
	minCode, maxCode := parseCodeFilter(code)

	// Populate table rows
	rowIndex := 1 // Start populating from the second row
	if url == "" {
		t.table.Clear()
		t.SetTableHeader(headers)
		for _, entry := range har.Log.Entries {
			if entry.Request.URL != url && entry.Response.Status >= minCode && entry.Response.Status <= maxCode {
				t.populateRow(rowIndex, entry)
				rowIndex++
			}
		}
	} else {
		t.table.Clear()
		t.SetTableHeader(headers)
		for _, entry := range har.Log.Entries {
			if strings.Contains(entry.Request.URL, url) && entry.Response.Status >= minCode && entry.Response.Status <= maxCode {
				t.populateRow(rowIndex, entry)
				rowIndex++
			}
		}
	}
	t.table.SetFocusFunc(func() {
		t.table.SetSelectable(true, false)
	})
	t.table.Select(1, 0).SetFixed(1, 1).SetSelectedFunc(
		func(row, column int) {
			entry := har.Log.Entries[row-1]
			t.showURLOptions(entry)
		})
}

func (t *Terminal) showURLOptions(entry Entry) {
	t.modal = tview.NewModal().SetText("Select an option").
		AddButtons([]string{"Request Headers", "Response Headers", "Content", "Timings", "Cancel"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			switch buttonLabel {
			case "Request Headers":
				t.showRequestDetails(entry)
			case "Response Headers":
				t.ShowResponseDetails(entry)
			case "Content":
				t.showContentDetails(entry)
			case "Timings":
				t.showTimingDetails(entry)
			case "Cancel":
				t.app.SetRoot(t.main, true)
			}
		})
	t.app.SetRoot(t.modal, true).EnableMouse(true)
}

func (t *Terminal) showRequestDetails(entry Entry) {
	requestView := tview.NewTextView().
		SetDynamicColors(true).
		SetText(fmt.Sprintf("[yellow]Request Headers:[white]\n%s", formatHeaders(entry.Request.Headers))).
		SetScrollable(true).
		SetWrap(true)
	requestView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEsc {
			t.app.SetRoot(t.modal, true)
		}
		return event
	})

	t.app.SetRoot(requestView, true).EnableMouse(true)
}

func (t *Terminal) ShowResponseDetails(entry Entry) {
	responseView := tview.NewTextView().
		SetDynamicColors(true).
		SetText(fmt.Sprintf("[yellow]Response Headers:[white]\n%s", formatHeaders(entry.Response.Headers))).
		SetScrollable(true).
		SetWrap(true)
	responseView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEsc {
			t.app.SetRoot(t.modal, true)
		}
		return event
	})

	t.app.SetRoot(responseView, true).EnableMouse(true)
}

func (t *Terminal) showTimingDetails(entry Entry) {
	timingsView := tview.NewTextView().
		SetDynamicColors(true).
		SetText(fmt.Sprintf("[yellow]Timings:[white]\n%s", formatTimings(entry.Timings))).
		SetScrollable(true).
		SetWrap(true)
	timingsView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEsc {
			t.app.SetRoot(t.modal, true)
		}
		return event
	})
	t.app.SetRoot(timingsView, true).EnableMouse(true)
}

func (t *Terminal) showContentDetails(entry Entry) {
	contentView := tview.NewTextView().
		SetDynamicColors(true).
		SetText(fmt.Sprintf("[yellow]Content:[white]\n%s", entry.Response.Content.Text)).
		SetScrollable(true).
		SetWrap(true)
	contentView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEsc {
			t.app.SetRoot(t.modal, true)
		}
		return event
	})
	t.app.SetRoot(contentView, true).EnableMouse(true)
}

func (t *Terminal) SetTableHeader(headers []string) {
	for col, header := range headers {
		t.table.SetCell(0, col, tview.NewTableCell(header).
			SetTextColor(tview.Styles.PrimaryTextColor).
			SetAlign(tview.AlignCenter).
			SetSelectable(false))
	}
}

func parseCodeFilter(code string) (int, int) {
	switch code {
	case "1XX":
		return 100, 199
	case "2XX":
		return 200, 299
	case "3XX":
		return 300, 399
	case "4XX":
		return 400, 499
	case "5XX":
		return 500, 599
	case "ALL":
		return 0, 999
	default:
		return 0, 999
	}
}

func (t *Terminal) populateRow(rowIndex int, entry Entry) {
	t.setTableCell(rowIndex, 0, fmt.Sprintf("%d", rowIndex), tview.AlignCenter, true)
	t.setTableCell(rowIndex, 1, entry.Request.Method, tview.AlignLeft, true)
	t.setTableCell(rowIndex, 2, fmt.Sprintf("%d", entry.Response.Status), tview.AlignCenter, true)
	t.setTableCell(rowIndex, 3, formatDomain(entry.Request.URL), tview.AlignLeft, true)
	t.setTableCell(rowIndex, 4, formatURL(entry.Request.URL), tview.AlignLeft, true)
	t.setTableCell(rowIndex, 5, entry.ServerIP, tview.AlignLeft, true)
	t.setTableCell(rowIndex, 6, entry.Connection, tview.AlignCenter, true)
	t.setTableCell(rowIndex, 7, fmt.Sprintf("%.2f", entry.Time), tview.AlignCenter, true)
}

// Helper function to set table cells
func (t *Terminal) setTableCell(row, col int, text string, align int, selectable bool) {
	t.table.SetCell(row, col, tview.NewTableCell(text).
		SetTextColor(tview.Styles.PrimaryTextColor).
		SetAlign(align).
		SetSelectable(selectable))
}

func (t *Terminal) CreateDropDown(har *HAR) {
	if har == nil || har.Log.Entries == nil {
		fmt.Println("Invalid HAR data")
		return
	}

	// Generate options based on unique response codes
	codeSet := map[string]struct{}{
		"ALL": {}, // Always include "ALL"
	}
	for _, entry := range har.Log.Entries {
		status := entry.Response.Status
		switch {
		case status >= 100 && status <= 199:
			codeSet["1XX"] = struct{}{}
		case status >= 200 && status <= 299:
			codeSet["2XX"] = struct{}{}
		case status >= 300 && status <= 399:
			codeSet["3XX"] = struct{}{}
		case status >= 400 && status <= 499:
			codeSet["4XX"] = struct{}{}
		case status >= 500 && status <= 599:
			codeSet["5XX"] = struct{}{}
		}
	}

	// Convert map keys to a sorted list
	var options []string
	for opt := range codeSet {
		options = append(options, opt)
	}
	sort.Strings(options)

	// Configure dropdown
	t.dropdown.SetLabel("Select an error code: ").
		SetOptions(options, func(option string, index int) {
			// When the user selects an option, update the table
			t.CreateTable(har, option, "")
		}).
		SetCurrentOption(0)
}

func (t *Terminal) CreateInputField(har *HAR) {
	if t.input == nil {
		t.input = tview.NewInputField()
	}
	t.input.SetLabel("Filter: ")
	t.input.SetFieldWidth(30)
	t.input.SetFieldTextColor(tview.Styles.PrimaryTextColor)

	t.input.SetChangedFunc(func(text string) {
		url := t.input.GetText()
		if url != "" {
			for _, entry := range har.Log.Entries {
				if strings.Contains(entry.Request.URL, url) {
					t.CreateTable(har, "", url)
					return
				}
			}
		} else {
			t.CreateTable(har, "", "")
		}
	})
}

func (t *Terminal) Layout() {
	primitives := []tview.Primitive{t.table, t.dropdown, t.input}
	form := tview.NewForm().AddFormItem(t.dropdown).AddFormItem(t.input)
	form.AddButton("Quit", func() {
		t.app.Stop()
	})
	form.SetBorder(true).SetTitle("Menu")
	t.main = tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(form, 60, 1, false).
		AddItem(t.table, 0, 1, true)

	t.main.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTab {
			for i, primitive := range primitives {
				if primitive == t.app.GetFocus() {
					t.app.SetFocus(primitives[(i+1)%len(primitives)])
					return nil
				}
			}
		}
		if event.Key() == tcell.KeyEsc {
			t.app.Stop()
			return nil
		}
		return event
	})

	t.app.SetRoot(t.main, true).EnableMouse(true)
}

func (t *Terminal) Run(har *HAR) error {
	// Ensure the terminal is initialized
	if t.app == nil || t.table == nil {
		return fmt.Errorf("terminal not initialized")
	}

	// Configure UI layout
	t.Layout()
	t.CreateInputField(har)
	t.CreateDropDown(har) // Configure the dropdown and link it to the table

	// Start the tview application loop
	if err := t.app.Run(); err != nil {
		return fmt.Errorf("application run failed: %w", err)
	}
	return nil
}
