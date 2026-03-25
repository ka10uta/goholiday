package goholiday

import "time"

const dateFormat = "2006-01-02"

type Schedule interface {
	GetNationalHolidays() map[string]string
	GetWeekdayHolidays() map[time.Weekday]struct{}
}

type Goholiday struct {
	schedule       Schedule
	uniqueHolidays map[string]struct{}
}

func New(schedule Schedule) *Goholiday {
	return &Goholiday{
		schedule:       schedule,
		uniqueHolidays: map[string]struct{}{},
	}
}

func (g *Goholiday) IsNationalHoliday(t time.Time) bool {
	_, exist := g.schedule.GetNationalHolidays()[t.Format(dateFormat)]
	return exist
}

func (g *Goholiday) IsHoliday(t time.Time) bool {
	if g.isWeekdayHoliday(t) {
		return true
	}
	dateKey := t.Format(dateFormat)
	if _, exist := g.schedule.GetNationalHolidays()[dateKey]; exist {
		return true
	}
	if _, exist := g.uniqueHolidays[dateKey]; exist {
		return true
	}
	return false
}

func (g *Goholiday) isWeekdayHoliday(t time.Time) bool {
	_, exist := g.schedule.GetWeekdayHolidays()[t.Weekday()]
	return exist
}

func (g *Goholiday) SetUniqueHolidays(ts []time.Time) {
	for _, t := range ts {
		g.uniqueHolidays[t.Format(dateFormat)] = struct{}{}
	}
}

func (g *Goholiday) IsBusinessDay(t time.Time) bool {
	return !g.IsHoliday(t)
}

func (g *Goholiday) BusinessDaysBefore(t time.Time, bds int) time.Time {
	return g.travelBusinessDays(t, bds, -1)
}

func (g *Goholiday) BusinessDaysAfter(t time.Time, bds int) time.Time {
	return g.travelBusinessDays(t, bds, 1)
}

func (g *Goholiday) travelBusinessDays(t time.Time, bds int, course int) time.Time {
	for tbds := 0; tbds != bds; {
		if t = t.AddDate(0, 0, course); !g.IsHoliday(t) {
			tbds++
		}
	}
	return t
}
