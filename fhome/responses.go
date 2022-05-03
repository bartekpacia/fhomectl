package fhome

import "fmt"

// Response is a websocket message sent from the server to the client in
// response to the client's previous websocket message to the server.

type Response struct {
	ActionName   string `json:"action_name"`
	RequestToken string `json:"request_token"`
	Status       string `json:"status"`
	Source       string `json:"source"`

	Details string `json:"details"` // Non-empty for "disconnecting" action
	Reason  string `json:"reason"`  // Non-empty for "disconnecting" action
}

type GetUserConfigResponse struct {
	ActionName   string `json:"action_name"`
	RequestToken string `json:"request_token"`
	Status       string `json:"status"`
	Source       string `json:"source"`
	File         string `json:"file"`
}

type GetMyResourcesResponse struct {
	ActionName    string `json:"action_name"`
	RequestToken  string `json:"request_token"`
	Status        string `json:"status"`
	Source        string `json:"source"`
	AvatarID0     string `json:"avatar_id_0"`
	FriendlyName0 string `json:"friendly_name_0"`
	ResourceType0 string `json:"resource_type_0"`
	UniqueID0     string `json:"unique_id_0"`
}

type TouchesResponse struct {
	ActionName string `json:"action_name"`
	Response   struct {
		ProjectVersion          string `json:"ProjectVersion"`
		Status                  bool   `json:"Status"`
		StatusText              string `json:"StatusText"`
		MobileDisplayProperties struct {
			Cells []MobileDisplayCell `json:"Cells"`
		} `json:"MobileDisplayProperties"`
	} `json:"response"`
	Status       string `json:"status"`
	Source       string `json:"source"`
	RequestToken string `json:"request_token"`
}

// MobileDisplayCell is a Cell, but returned from "touches" action.
type MobileDisplayCell struct {
	// Cell description. Note that this is by the configurator app, not by the
	// user in the mobile or web app.
	Desc string `json:"CD"`
	// Object ID
	ID string `json:"OI"`
	// Type number. Known values: 706, 707, 708, 709, 710, 711, 717, 718, 719,
	// 722, 724, 760
	TypeNumber string `json:"TN"`
	// Preset. Known values: 0, 1, 4
	Preset string `json:"P"`
	// Style. Display Type TEMP always has this set to 2.
	Style string `json:"Se"`
	// Minimum value
	MinVal string `json:"Min"`
	// Maximum value
	MaxVal string `json:"Max"`
	// Step (aka current value). Display Type TEMP always has this set to
	// 0xa005.
	Step string `json:"Sp"`
	// Display Type. Known values: BIT, BYTE, TEMP (Temperature), PROC
	// (Percentage), RGB (Light)
	DisplayType string `json:"DT"`
	// Cell permission. Known values: FC (Full Control), RO (Read Only)
	Permission string `json:"CP"`
}

func (cell MobileDisplayCell) String() string {
	return fmt.Sprintf("id: %s, desc: %s, type: %s, preset: %s, style: %s, perm: %s, step/value: %s\n",
		cell.ID, cell.Desc, cell.DisplayType, cell.Preset, cell.Style, cell.Permission, cell.Step,
	)
}

type StatusTouchesChangedResponse struct {
	ActionName string `json:"action_name"`
	Response   struct {
		ProjectVersion string      `json:"ProjectVersion"`
		Status         bool        `json:"Status"`
		StatusText     string      `json:"StatusText"`
		CellValues     []CellValue `json:"CV"`
		ServerTime     int         `json:"ServerTime"`
	} `json:"response"`
	Status string `json:"status"`
	Source string `json:"source"`
}

type CellValue struct {
	ID  string `json:"VOI"`
	Ii  string `json:"II"`
	Dt  string `json:"DT"`
	Dv  string `json:"DV"`
	Dvs string `json:"DVS"`
}
