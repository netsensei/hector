package ui

type Tab struct {
	URL     string
	Status  string
	Content string
	// History History
	// State string
}

type Tabs struct {
	Tabs      []Tab
	ActiveTab int
}

func NewTabs() *Tabs {
	return &Tabs{
		ActiveTab: 0,
	}
}

func (ts *Tabs) Add(tab Tab) {
	if len(ts.Tabs) == 0 {
		ts.Tabs = append(ts.Tabs, tab)
	} else {
		if ts.ActiveTab == len(ts.Tabs)-1 {
			ts.Tabs = append(ts.Tabs, tab)
		} else {
			tabs := append(ts.Tabs, Tab{})
			copy(tabs[ts.ActiveTab+1:], tabs[ts.ActiveTab:])
			tabs[ts.ActiveTab+1] = tab
			ts.Tabs = tabs
		}
		ts.ActiveTab++
	}
}

func (ts *Tabs) Remove() {
	if len(ts.Tabs) > 1 {
		ts.Tabs = append(ts.Tabs[:ts.ActiveTab], ts.Tabs[ts.ActiveTab+1:]...)
		ts.ActiveTab--
	}
}

func (ts *Tabs) Up() {
	ts.ActiveTab++
	if ts.ActiveTab >= len(ts.Tabs)-1 {
		ts.ActiveTab = len(ts.Tabs) - 1
	}
}

func (ts *Tabs) Down() {
	if ts.ActiveTab != 0 {
		ts.ActiveTab--
	}
}

func (ts *Tabs) Current() (*Tab, int) {
	return &ts.Tabs[ts.ActiveTab], ts.ActiveTab
}

func (ts *Tabs) Update(tab Tab) {
	ts.Tabs[ts.ActiveTab] = tab
}

func (ts *Tabs) Count() int {
	return len(ts.Tabs)
}
