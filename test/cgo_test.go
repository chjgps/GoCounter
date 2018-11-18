package main

import (
	"encoding/json"
	"fmt"
	"testing"
)

// 关门传感器结果返回

const testJson =

/*`
	{
		"door_id":1,
		"weight":[
			{"ch":0,"value":0},{"ch":1,"value":0}
		],
		"ultra_sonic":[
			{"ch":0,"value":0},{"ch":1,"value":0}
		],
		"image":[
			{"type":0,"count":0},{"type":1,"count":0}
		],
		"rfid":[
			"E2003069911502140910BE6A",
			"201311010000000000000000"
		]
	}
`
*/

`{"door_id":1,"ultra_sonic":[],"weight":[{"ch":0,"value":1204}],"image":[],"rfid":[]}`

type SensorData struct {
	DoorId      int           `json:"door_id"`
	Weights     []*Weight     `json:"weight"`
	UltraSonics []*UltraSonic `json:"ultra_sonic"`
	Images      []*Image      `json:"image"`
	Rfids       []string      `json:"rfid"`
}

func (s *SensorData) String() string {
	bytes, err := json.MarshalIndent(s, "  ", "    ")
	if err != nil {
		return fmt.Sprintf("%#v", err)
	}

	return string(bytes)
}

type Weight struct {
	Channel int `json:"ch"`
	Value   int `json:"value"`
}

type UltraSonic struct {
	Channel int `json:"ch"`
	Value   int `json:"value"`
}

type Image struct {
	Type  int `json:"type"`
	Count int `json:"count"`
}

func TestJson(t *testing.T) {
	obj := &SensorData{}
	if err := json.Unmarshal([]byte(testJson), obj); err != nil {
		t.Fatal(err)
		return
	}

	t.Log(obj.String())
}
