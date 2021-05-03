package Watch

import (
	"fmt"
	"time"
)

type Calendar struct {
	Year          int
	Month         time.Month
	Day           int
	Today         int
	Week          int
	Complex       time.Time
	MonthDays     []time.Time
	TodayRow      int
	CalendarRows  [][]string
	CalendarTitle string
}

func NewCalendar(params ...int) Calendar {
	var now time.Time
	if len(params) > 6 {
		now = time.Date(
			params[0],
			time.Month(params[1]),
			params[2],
			params[3],
			params[4],
			params[5],
			params[6],
			time.Local,
		)
	} else {
		now = time.Now()
	}
	year, month, today := now.Date()
	_, week := now.ISOWeek()
	days := []time.Time{}
	for i := 1; i < 32; i++ {
		switch month.String() {
		case "January", "March", "May", "July", "August", "October", "December":
			days = append(days, time.Date(year, month, i, 1, 0, 0, 0, time.Local))
		case "February":
			if i < 29 {
				days = append(days, time.Date(year, month, i, 1, 0, 0, 0, time.Local))
			} else if i == 29 {
				if year%4 == 0 {
					if year%100 == 0 {
						if year%400 == 0 {
							days = append(days, time.Date(year, month, i, 1, 0, 0, 0, time.Local))
						} else {
							continue
						}
					} else {
						days = append(days, time.Date(year, month, i, 1, 0, 0, 0, time.Local))
					}
				}
			}
		case "April", "June", "September", "November":
			if i < 31 {
				days = append(days, time.Date(year, month, i, 1, 0, 0, 0, time.Local))
			}
		}
	}

	calendarRows := [][]string{
		{"Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"},
	}
	var rowEmphasize int
	calDays := make([]string, 7)
	for i, day := range days {
		switch day.Weekday().String() {
		case "Sunday":
			calDays[0] = fmt.Sprint(int(day.Day()))
		case "Monday":
			calDays[1] = fmt.Sprint(int(day.Day()))
		case "Tuesday":
			calDays[2] = fmt.Sprint(int(day.Day()))
		case "Wednesday":
			calDays[3] = fmt.Sprint(int(day.Day()))
		case "Thursday":
			calDays[4] = fmt.Sprint(int(day.Day()))
		case "Friday":
			calDays[5] = fmt.Sprint(int(day.Day()))
		case "Saturday":
			calDays[6] = fmt.Sprint(int(day.Day()))
		}
		if len(days)-1 == i || day.Weekday().String() == "Saturday" {
			calendarRows = append(calendarRows, calDays)
			calDays = make([]string, 7)
		}
		// if today == day.Day() {
		// 	rowEmphasize = len(calendarRows)
	}

Loop:
	for i, v := range calendarRows {
		for _, d := range v {
			if d == fmt.Sprint(today) {
				rowEmphasize = i
				break Loop
			}
		}
	}

	return Calendar{
		Year:          year,
		Month:         month,
		Day:           today,
		Today:         time.Now().Day(),
		Week:          week,
		MonthDays:     days,
		Complex:       now,
		TodayRow:      rowEmphasize,
		CalendarRows:  calendarRows,
		CalendarTitle: month.String() + " - " + fmt.Sprint(year) + " - " + fmt.Sprint(week) + " week",
	}
}
