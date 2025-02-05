package models

type User struct {
	CurrentUser struct {
		Username       string `json:"username"`
		PrivilegeLevel int    `json:"privilegeLevel"`
	} `json:"currentUser"`
	CurrentProjectName string            `json:"currentProjectName"`
	Language           string            `json:"language"`
	ProjectOverviews   []ProjectOverview `json:"projectOverviews"`
	Projects           []Project         `json:"projects"`
}

type ProjectOverview struct {
	Name         string `json:"name"`
	Instances    int    `json:"instances"`
	Sessions     int    `json:"sessions"`
	Health       string `json:"health"`
	Streak       string `json:"streak"`
	Pause        string `json:"pause"`
	LastActivity string `json:"lastActivity"`
}

type Project struct {
	Name         string       `json:"name"`
	Health       string       `json:"health"`
	Streak       string       `json:"streak"`
	LastActivity string       `json:"lastActivity"`
	Pause        string       `json:"pause"`
	Ingredients  []Ingredient `json:"ingredients"`
	Sessions     []Session    `json:"sessions"`
	Own          bool         `json:"own"`
}

type Ingredient struct {
	Name      string     `json:"name"`
	Active    bool       `json:"active"`
	Instances []Instance `json:"instances"`
}

type Instance struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Pos           Pos    `json:"pos"`
	Active        bool   `json:"active"`
	SessionNumber int    `json:"sessionNumber"`
}

type Session struct {
	ID                  string               `json:"id"`
	SessionNumber       int                  `json:"sessionNumber"`
	Active              bool                 `json:"active"`
	Order               int                  `json:"order"`
	SessionTitle        string               `json:"sessionTitle"`
	StartDate           string               `json:"startDate"`
	EndDate             string               `json:"endDate"`
	Duration            string               `json:"duration"`
	CumulativeDuration  string               `json:"cumulativeDuration"`
	Description         string               `json:"description"`
	AssociatedInstances []AssociatedInstance `json:"associatedInstances"`
	RestartDate         string               `json:"restartDate"`
}

type AssociatedInstance struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Pos           Pos    `json:"pos"`
	Active        bool   `json:"active"`
	SessionNumber int    `json:"sessionNumber"`
}

type Pos struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}
