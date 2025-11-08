package time

import (
    "database/sql/driver"
    "encoding/json"
    "github.com/pkg/errors"
    goTime "time"
    "strconv"
    "time"
)

type Time time.Time

func (t *Time) UnmarshalJSON(data []byte) error {
    millis, err := strconv.ParseInt(string(data), 10, 64)
    if err != nil {
        return err
    }
    *t = Time(ConvertInt642Time(millis))
    return nil
}

func (t *Time) MarshalJSON() ([]byte, error) {
    millis := ConvertTime2Int64(time.Time(*t))
    return json.Marshal(millis)
}

func (t Time) Value() (driver.Value, error) {
    return time.Time(t), nil
}

func (t Time) Format() string {
    return time.Time(t).Format(time.RFC3339)
}

func (t Time) FormatShortly() string {
    return time.Time(t).Format("2006-01-02 15:04:05")
}

func (t *Time) Scan(value interface{}) error {
    if v, ok := value.(time.Time); ok {
        *t = Time(v)
        return nil
    } else {
        return errors.Errorf("type convert error : %+v", value)
    }
}

func (t Time) After(u Time) bool {
    a := time.Time(t)
    b := time.Time(u)
    return a.After(b)
}

func (t Time) Transfer() goTime.Time {
    return goTime.Time(t)
}

//ConvertTime2Int64 is used to convert time.Time to millisecond since Jan 01 1970. (UTC)
//Example: https://play.golang.org/p/Aa8eg_rdGLF
func ConvertTime2Int64(t time.Time) int64 {
    return t.UnixNano() / int64(time.Millisecond)
}

//ConvertInt642Time is used to convert millisecond since Jan 01 1970. (UTC) to time.Time
//Example: https://play.golang.org/p/Aa8eg_rdGLF
func ConvertInt642Time(i int64) time.Time {
    msIn1Second := int64(10 * 10 * 10)

    sec := i / msIn1Second
    nsec := (i % msIn1Second) * int64(time.Millisecond)

    result := time.Unix(sec, nsec)

    return result
}
