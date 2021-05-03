package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	Watch "github.com/buahaha/watch/Watch"
	"github.com/gen2brain/beeep"
	"github.com/recoilme/slowpoke"

	ui "github.com/buahaha/termui/v3"
	"github.com/buahaha/termui/v3/widgets"
)

func main() {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	swDisplay := widgets.NewParagraph()
	stopwatch := Watch.NewStopwatch()
	swDisplay.Text = stopwatch.Diff().String() //time.Duration(4e17).String()
	swDisplay.Title = "Time elapsed"
	swDisplay.SetRect(0, 0, 21, 3)

	swStart := widgets.NewParagraph()
	swStart.Text = "[Start!](mod:bold)"
	swStart.SetRect(0, 3, 8, 6)

	swClear := widgets.NewParagraph()
	swClear.Text = "[Sum & Clear](mod:underline)"
	swClear.SetRect(8, 3, 21, 6)

	l := widgets.NewList()
	l.Title = "Stopwatch laps"
	l.SetRect(21, 0, 46, 12)
	l.Rows = []string{}

	swBars := widgets.NewBarChart()
	swBars.Title = "Laps difference"
	swBars.SetRect(46, 0, 80, 12)
	swBars.Data = []float64{}
	swBars.Labels = []string{}

	localTimeDisplay := widgets.NewParagraph()
	localTimeDisplay.SetRect(0, 6, 21, 9)
	localTimeDisplay.Title = "Local time"
	// localTimeHour, localTimeMinute, localTimeSecond := Watch.LocalTime().Clock()
	localTimeDisplay.Text = Watch.LocalTime().String()[:19] //fmt.Sprint(localTimeHour, ":", localTimeMinute, ":", localTimeSecond)

	utcTimeDisplay := widgets.NewParagraph()
	utcTimeDisplay.SetRect(0, 9, 21, 12)
	utcTimeDisplay.Title = "UTC/GMT time"
	// utcTimeHour, utcTimeMinute, utcTimeSecond := Watch.UniversalTime().Clock()
	utcTimeDisplay.Text = Watch.UniversalTime().String()[:19] //fmt.Sprint(utcTimeHour, ":", utcTimeMinute, ":", utcTimeSecond)

	timer := Watch.NewTimer()

	timerDisplay := widgets.NewParagraph()
	timerDisplay.SetRect(0, 12, 21, 15)
	timerDisplay.Title = "Timer"
	timerDisplay.Text = "Set below."

	timerSetHours := widgets.NewInput()
	timerSetHours.SetRect(0, 15, 5, 18)
	timerSetHours.Title = "H"
	timerSetHours.TitleStyle.Fg = ui.ColorRed
	timerSetHours.Text = ""
	// var timerHours int

	timerSetMinutes := widgets.NewInput()
	timerSetMinutes.SetRect(5, 15, 12, 18)
	timerSetMinutes.Title = "Min"
	timerSetMinutes.TitleStyle.Fg = ui.ColorRed
	timerSetMinutes.Text = ""
	// var timerMinutes int

	timerSetSeconds := widgets.NewInput()
	timerSetSeconds.SetRect(12, 15, 21, 18)
	timerSetSeconds.Title = "Sec"
	timerSetSeconds.TitleStyle.Fg = ui.ColorRed
	timerSetSeconds.Text = ""
	// var timerSeconds int

	beeep.DefaultFreq = 432.0

	input := widgets.NewInput()
	// input.WrapText = true
	input.SetRect(0, 21, 21, 25)
	input.Title = "insert msg here"

	timerStart := widgets.NewParagraph()
	timerStart.SetRect(0, 18, 21, 21)
	timerStart.Text = "Run the timer..."
	timerStart.Title = "Set above"

	cal := Watch.NewCalendar()
	calendar := widgets.NewTable()
	calendar.Title = cal.CalendarTitle
	calendar.SetRect(21, 12, 54, 21)

	calendar.Rows = cal.CalendarRows

	calendar.TextAlignment = ui.AlignCenter
	calendar.BorderStyle.Fg = ui.ColorYellow
	calendar.RowSeparator = false
	calendar.RowStyles[cal.TodayRow] = ui.NewStyle(ui.ColorGreen, ui.ColorClear, ui.ModifierBold)
	calendar.ColumnWidths = []int{3, 3, 3, 3, 3, 3, 3}
	calendar.ColumnWidths[cal.Complex.Weekday()] = 7

	goToDate := widgets.NewInput()
	goToDate.SetRect(21, 21, 54, 25)
	goToDate.Title = "Year Month Day exchange field"

	dock := widgets.NewInput()
	dock.SetRect(54, 12, 80, 25)
	dock.Title = fmt.Sprintf("Note %d-%d-%d", cal.Year, cal.Month, cal.Day)
	// create database
	file := "db/Watch.db"
	// close all opened database
	defer slowpoke.CloseAll()
	res, _ := slowpoke.Get(file, []byte(strings.Split(dock.Title, " ")[1]))
	dock.Text = string(res)

	ui.Render(
		swDisplay, swStart, swClear, l, swBars,
		localTimeDisplay, utcTimeDisplay,
		timerDisplay, timerSetHours, timerSetMinutes, timerSetSeconds, timerStart, input,
		calendar, goToDate, //yearBack, yearForward, monthBack, monthForward, thisMonth,
		dock,
	)

	inputFields := []widgets.Input{
		*timerSetHours,
		*timerSetMinutes,
		*timerSetSeconds,
		*input,
		*dock,
	}

	uiEvents := ui.PollEvents()

	swTicker := time.NewTicker(100 * time.Millisecond)
	swTickerChan := swTicker.C
	swTicker.Stop()
	clockTicker := time.NewTicker(time.Second).C
	calendarTicker := time.NewTicker(5 * time.Second).C
	// deathToll := len(myDeaths)
	renderTicker := time.NewTicker(1000 / 23 * time.Millisecond)
	for {
		select {
		case e := <-uiEvents:
			switch e.ID { // event string/identifier
			case "q", "<C-c>": // press 'q' or 'C-c' to quit
				return
			case "<MouseLeft>":
				for _, i := range inputFields {
					if i.Focusing {
						i.NoFocus()
					}
				}
				payload := e.Payload.(ui.Mouse)
				x, y := payload.X, payload.Y
				if x >= 8 && x <= 21 &&
					y >= 3 && y <= 6 {
					var duration time.Duration
					lastLap := stopwatch.Stop().String()
					swStart.Text = "[Start!](mod:bold)"
					l.Rows = append(l.Rows, fmt.Sprint(len(l.Rows)+1, ". ", lastLap))
					stopwatch = Watch.NewStopwatch()
					for _, s := range l.Rows {
						intermediateDuration, _ := time.ParseDuration(strings.Split(s, " ")[1])
						duration = intermediateDuration + duration
					}
					swDisplay.Text = duration.String()
					// ui.Render(swStart, swDisplay)
					l.Rows = []string{}
					swBars.Data = []float64{}
					swBars.Labels = []string{}
					// ui.Render(l, swBars)
				} else if x >= 0 && x <= 8 &&
					y >= 0 && y <= 6 &&
					!stopwatch.Running ||
					x >= 0 && x <= 21 &&
						y >= 0 && y <= 3 &&
						!stopwatch.Running {
					stopwatch.Start()
					swTicker.Reset(50 * time.Millisecond)
					swStart.Text = "[Stop!](mod:bold)"
					// ui.Render(swStart)
				} else if x >= 0 && x <= 8 &&
					y >= 3 && y <= 6 &&
					stopwatch.Running ||
					x >= 0 && x <= 21 &&
						y >= 0 && y <= 3 &&
						stopwatch.Running {
					swDisplay.Text = stopwatch.Stop().String()
					swTicker.Stop()
					stopwatch = Watch.NewStopwatch()
					ui.Render(swDisplay)
					swStart.Text = "[Start!](mod:bold)"
					ui.Render(swStart)
					l.Rows = append(l.Rows, fmt.Sprint(len(l.Rows)+1, ". ", swDisplay.Text))
					// ui.Render(l)
					bar, _ := time.ParseDuration(swDisplay.Text)
					swBars.Data = append(swBars.Data, float64(bar))
					swBars.Labels = append(swBars.Labels, fmt.Sprint(len(l.Rows)))
					// ui.Render(swBars)
				} else if x >= 21 && x <= 42 &&
					y >= 6 && y <= 12 {
					l.ScrollHalfPageDown()
					// ui.Render(l)
				} else if x >= 21 && x <= 42 &&
					y >= 0 && y <= 5 {
					l.ScrollHalfPageUp()
					// ui.Render(l)
				} else if x >= 0 && x <= 5 &&
					y >= 15 && y <= 18 {
					// go timerSetMinutes.NoFocus()
					_, err := strconv.Atoi(timerSetMinutes.Text)
					if timerSetMinutes.Text == "" || err != nil {
						timerSetMinutes.TitleStyle.Fg = ui.ColorRed
					} else {
						timerSetMinutes.TitleStyle.Fg = ui.ColorGreen
						timerStart.Title = "Ready"
						timerStart.TitleStyle.Fg = ui.ColorGreen
					}
					// go timerSetSeconds.NoFocus()
					_, err = strconv.Atoi(timerSetSeconds.Text)
					if timerSetSeconds.Text == "" || err != nil {
						timerSetSeconds.TitleStyle.Fg = ui.ColorRed
					} else {
						timerSetSeconds.TitleStyle.Fg = ui.ColorGreen
						timerStart.Title = "Ready"
						timerStart.TitleStyle.Fg = ui.ColorGreen
					}
					// go input.NoFocus()
					timerSetHours.Text = ""
					go timerSetHours.Focus()
					_, err = strconv.Atoi(timerSetHours.Text)
					if err != nil {
						timerSetHours.Text = ""
						timerSetHours.TitleStyle.Fg = ui.ColorRed
					} else {
						timerSetHours.TitleStyle.Fg = ui.ColorGreen
						timerStart.Title = "Ready"
						timerStart.TitleStyle.Fg = ui.ColorGreen
					}
					// ui.Render(timerSetHours, timerSetMinutes, timerSetSeconds, timerStart)
					// timerSetHours.NoFocus()
				} else if x >= 5 && x <= 12 &&
					y >= 15 && y <= 18 {
					// go timerSetHours.NoFocus()
					_, err := strconv.Atoi(timerSetHours.Text)
					if timerSetHours.Text == "" || err != nil {
						timerSetHours.TitleStyle.Fg = ui.ColorRed
					} else {
						timerSetHours.TitleStyle.Fg = ui.ColorGreen
						timerStart.Title = "Setting done"
						timerStart.TitleStyle.Fg = ui.ColorGreen
					}
					// go timerSetSeconds.NoFocus()
					_, err = strconv.Atoi(timerSetSeconds.Text)
					if timerSetSeconds.Text == "" || err != nil {
						timerSetSeconds.TitleStyle.Fg = ui.ColorRed
					} else {
						timerSetSeconds.TitleStyle.Fg = ui.ColorGreen
						timerStart.Title = "Setting done"
						timerStart.TitleStyle.Fg = ui.ColorGreen
					}
					// go input.NoFocus()
					timerSetMinutes.Text = ""
					go timerSetMinutes.Focus()
					_, err = strconv.Atoi(timerSetMinutes.Text)
					if err != nil {
						timerSetMinutes.Text = ""
						timerSetMinutes.TitleStyle.Fg = ui.ColorRed
					} else {
						timerSetMinutes.TitleStyle.Fg = ui.ColorGreen
						timerStart.Title = "Ready"
						timerStart.TitleStyle.Fg = ui.ColorGreen
					}
					// ui.Render(timerSetHours, timerSetMinutes, timerSetSeconds, timerStart)
					// go timerSetMinutes.NoFocus()
				} else if x >= 12 && x <= 21 &&
					y >= 15 && y <= 18 {
					// go timerSetHours.NoFocus()
					_, err := strconv.Atoi(timerSetHours.Text)
					if err != nil {
						timerSetHours.Text = ""
						timerSetHours.TitleStyle.Fg = ui.ColorRed
					} else {
						timerSetHours.TitleStyle.Fg = ui.ColorGreen
						timerStart.Title = "Ready"
						timerStart.TitleStyle.Fg = ui.ColorGreen
					}
					// go timerSetMinutes.NoFocus()
					_, err = strconv.Atoi(timerSetMinutes.Text)
					if err != nil {
						timerSetMinutes.Text = ""
						timerSetMinutes.TitleStyle.Fg = ui.ColorRed
					} else {
						timerSetMinutes.TitleStyle.Fg = ui.ColorGreen
						timerStart.Title = "Ready"
						timerStart.TitleStyle.Fg = ui.ColorGreen
					}
					timerSetMinutes.TitleStyle.Fg = ui.ColorRed
					// go input.NoFocus()
					timerSetSeconds.Text = ""
					go timerSetSeconds.Focus()
					_, err = strconv.Atoi(timerSetSeconds.Text)
					if err != nil {
						timerSetSeconds.Text = ""
						timerSetSeconds.TitleStyle.Fg = ui.ColorRed
					} else {
						timerSetSeconds.TitleStyle.Fg = ui.ColorGreen
						timerStart.Title = "Ready"
						timerStart.TitleStyle.Fg = ui.ColorGreen
					}
					// ui.Render(timerSetHours, timerSetMinutes, timerSetSeconds, timerStart)
					// go timerSetSeconds.NoFocus()
				} else if x >= 0 && x <= 21 &&
					y >= 18 && y <= 21 {
					if !timer.Running {
						h, _ := strconv.Atoi(timerSetHours.Text)
						m, _ := strconv.Atoi(timerSetMinutes.Text)
						s, _ := strconv.Atoi(timerSetSeconds.Text)
						if h != 0 || m != 0 || s != 0 {
							timer.SetEndTime(h, m, s)
						}
						timerDisplay.Text = "Starting..."
						timer.Running = true
						timerDisplay.TitleStyle.Fg = ui.ColorGreen
						// ui.Render(timerDisplay)
						timerSetHours.TitleStyle.Fg = ui.ColorRed
						timerSetHours.Text = ""
						timerSetMinutes.TitleStyle.Fg = ui.ColorRed
						timerSetMinutes.Text = ""
						timerSetSeconds.TitleStyle.Fg = ui.ColorRed
						timerSetSeconds.Text = ""
						// ui.Render(timerSetHours, timerSetMinutes, timerSetSeconds)
						timerStart.Text = "Stop countdown"
						timerStart.TitleStyle.Fg = ui.ColorClear
						timerStart.Title = "Alert onGOing"
						// ui.Render(timerStart)
					} else if timer.Running {
						timerStart.Text = "Run the timer..."
						timerDisplay.Text = "Set again."
						timerDisplay.TitleStyle.Fg = ui.ColorClear
						timer.Stop()
						timer = Watch.NewTimer()
						// ui.Render(timerDisplay, timerStart)
					}
				} else if x >= 0 && x <= 21 &&
					y >= 21 && y <= 24 {
					// go timerSetHours.NoFocus()
					// go timerSetMinutes.NoFocus()
					// go timerSetSeconds.NoFocus()
					// go dock.NoFocus()
					go input.Focus()
					// go input.NoFocus()
				} else if x >= 54 && y >= 12 && x <= 80 && y <= 25 {
					go dock.Focus()
					// go dock.NoFocus()
				} else if x >= 21 && y >= 21 && x <= 54 && y <= 25 {
					go goToDate.Focus()
					// go goToDate.NoFocus()

				}
			}
			// case "<Resize>":
			// 	payload := e.Payload.(ui.Resize)
			// 	width, height := payload.Width, payload.Height
			// switch e.Type {
			// case ui.KeyboardEvent: // handle all key presses
			// 	// eventID := e.ID // keypress string
			// }
		// use Go's built-in tickers for updating and drawing data
		case <-renderTicker.C:
			ui.Render(
				swDisplay, swStart, swClear, l, swBars,
				localTimeDisplay, utcTimeDisplay,
				timerDisplay, timerSetHours, timerSetMinutes, timerSetSeconds, timerStart, input,
				calendar, goToDate, //yearBack, yearForward, monthBack, monthForward, thisMonth,
				dock,
			)
		case <-calendarTicker:
			key := []byte(strings.Split(dock.Title, " ")[1])
			val := []byte(dock.Text)
			//store
			slowpoke.Set(file, key, val)
			if time.Now().Day() != cal.Today {
				cal = Watch.NewCalendar(time.Now().Year(), int(time.Now().Month()), time.Now().Day(), 0, 0, 0, 0)
				calendar = widgets.NewTable()
				calendar.Title = cal.CalendarTitle
				calendar.SetRect(21, 12, 54, 21)

				calendar.Rows = cal.CalendarRows

				calendar.TextAlignment = ui.AlignCenter
				calendar.BorderStyle.Fg = ui.ColorYellow
				calendar.RowSeparator = false
				calendar.RowStyles[cal.TodayRow] = ui.NewStyle(ui.ColorGreen, ui.ColorClear, ui.ModifierBold)
				calendar.ColumnWidths = []int{3, 3, 3, 3, 3, 3, 3}
				calendar.ColumnWidths[cal.Complex.Weekday()] = 7

				dock.Title = fmt.Sprintf("Note %d-%d-%d", cal.Year, cal.Month, cal.Day)
				// get
				res, _ := slowpoke.Get(file, []byte(strings.Split(dock.Title, " ")[1]))
				dock.Text = string(res)
			}
			if goToDate.Text != "" {
				strArr := strings.Split(goToDate.Text, " ")
				var intArr []int
				for _, str := range strArr {
					tmp, err := strconv.Atoi(str)
					if err != nil {
						break
					}
					intArr = append(intArr, tmp)
				}
				if len(intArr) == 3 {
					cal = Watch.NewCalendar(intArr[0], intArr[1], intArr[2], 0, 0, 0, 0)
					calendar = widgets.NewTable()
					calendar.Title = cal.CalendarTitle
					calendar.SetRect(21, 12, 54, 21)

					calendar.Rows = cal.CalendarRows

					calendar.TextAlignment = ui.AlignCenter
					calendar.BorderStyle.Fg = ui.ColorYellow
					calendar.RowSeparator = false
					calendar.RowStyles[cal.TodayRow] = ui.NewStyle(ui.ColorGreen, ui.ColorClear, ui.ModifierBold)
					calendar.ColumnWidths = []int{3, 3, 3, 3, 3, 3, 3}
					calendar.ColumnWidths[cal.Complex.Weekday()] = 7
					// ui.Render(calendar)
					goToDate.Text = ""
					dock.Title = fmt.Sprintf("Note %d-%d-%d", cal.Year, cal.Month, cal.Day)
					// get
					res, _ := slowpoke.Get(file, []byte(strings.Split(dock.Title, " ")[1]))
					dock.Text = string(res)
				}
			}
		case <-swTickerChan:
			if stopwatch.Running {
				swDisplay.Text = stopwatch.Diff().String()
				ui.Render(swDisplay)
			}
		case <-clockTicker:
			// localTimeHour, localTimeMinute, localTimeSecond := Watch.LocalTime().Clock()
			localTimeDisplay.Text = Watch.LocalTime().String()[:19] //fmt.Sprint(localTimeHour, ":", localTimeMinute, ":", localTimeSecond)
			// utcTimeHour, utcTimeMinute, utcTimeSecond := Watch.UniversalTime().Clock()
			utcTimeDisplay.Text = Watch.UniversalTime().String()[:19] //fmt.Sprint(utcTimeHour, ":", utcTimeMinute, ":", utcTimeSecond)
			ui.Render(localTimeDisplay, utcTimeDisplay)
			if timer.Ended && timer.Running && time.Until(timer.EndTime) > -(5*time.Second) {
				timerDisplay.Text = "The end"
				timerDisplay.TitleStyle.Fg = ui.ColorRed
				ui.Render(timerDisplay)
				err := beeep.Alert("Alarm", input.Text, "")
				// time.Sleep(500 * time.Millisecond)
				if err != nil {
					log.Fatal(err)
				}
				beeep.Beep(432.0, 200)
				if err != nil {
					log.Fatal(err)
				}
				timerStart.Text = "Run the timer..."
				ui.Render(timerStart, input)
				timer.Stop()
			} else if timer.Running && time.Until(timer.EndTime) > -(5*time.Second) {
				timerDisplay.Text = timer.Countdown()
				ui.Render(timerDisplay)
				break
			} else {
				timerDisplay.Text = "The end"
				timerDisplay.TitleStyle.Fg = ui.ColorRed
				ui.Render(timerDisplay)
				timer = Watch.NewTimer()
				break
			}
		}
	}
}
