package timex

import (
	"time"

	"github.com/beego/ms304w-client/basis/errors"
)

// ----------------------
// FormatDate

func FormatTime(date string) (time.Time, error) {
	t, err := time.ParseInLocation("2006-01-02 15:04:05", date, time.Local)
	if err != nil {
		return time.Time{}, err
	}

	return t, nil
}

func String() string {
	t := time.Now()

	// return t.Add(time.Duration(8) * time.Hour).Format("2006-01-02 15:04:05")
	return t.Format("2006-01-02 15:04:05")
}

// ----------------------
// Time

type Time struct {
	time.Time
}

func (t Time) MarshalJSON() ([]byte, error) {
	if y := t.Year(); y < 0 || y >= 10000 {
		return nil, errors.New("Time.MarshalJSON: year outside of range [0,9999]")
	}

	if t.Unix() == 0 {
		return []byte(`"` + "" + `"`), nil
	}

	return []byte(t.Format(`"` + "2006-01-02 15:04:05" + `"`)), nil
}

func (t Time) String() string {
	return t.Format("2006-01-02 15:04:05")
}

// -----------------------
// Date

type Date struct {
	time.Time
}

func (d Date) MarshalJSON() ([]byte, error) {
	if y := d.Year(); y < 0 || y >= 10000 {
		return nil, errors.New("Time.MarshalJSON: year outside of range [0,9999]")
	}
	return []byte(d.Format(`"` + "2006-01-02" + `"`)), nil
}

func (d Date) String() string {
	return d.Format("2006-01-02")
}

// ----------------------
// Hour

type Hour struct {
	time.Time
}

func (h Hour) MarshalJSON() ([]byte, error) {
	if y := h.Year(); y < 0 || y >= 10000 {
		return nil, errors.New("Time.MarshalJSON: year outside of range [0,9999]")
	}
	return []byte(h.Format(`"` + "15:04:05" + `"`)), nil
}

func (h Hour) String() string {
	return h.Format("15:04:05")
}
