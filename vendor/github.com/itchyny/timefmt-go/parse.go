package timefmt

import (
	"errors"
	"fmt"
	"strconv"
	"time"
)

type parseError struct {
	source, format string
	err            error
}

func (err *parseError) Error() string {
	return fmt.Sprintf("failed to parse %q with %q: %s", err.source, err.format, err.err)
}

// Parse time string using the format.
func Parse(source, format string) (t time.Time, err error) {
	year, month, day, hour, min, sec, nsec, loc := 1900, 1, 1, 0, 0, 0, 0, time.UTC
	defer func() {
		if err != nil {
			err = &parseError{source, format, err}
		}
	}()
	var j, diff, century, yday int
	var pm bool
	var pending string
	for i, l := 0, len(source); i < len(format); i++ {
		if b := format[i]; b == '%' {
			i++
			if i == len(format) {
				err = errors.New("stray %")
				return
			}
			b = format[i]
		L:
			switch b {
			case 'Y':
				if year, diff, err = parseNumber(source[j:], 4, 'Y'); err != nil {
					return
				}
				j += diff
			case 'y':
				if year, diff, err = parseNumber(source[j:], 2, 'y'); err != nil {
					return
				}
				j += diff
				if year < 69 {
					year += 2000
				} else {
					year += 1900
				}
			case 'C':
				if century, diff, err = parseNumber(source[j:], 2, 'C'); err != nil {
					return
				}
				j += diff
			case 'g':
				if year, diff, err = parseNumber(source[j:], 2, b); err != nil {
					return
				}
				j += diff
				year += 2000
			case 'G':
				if year, diff, err = parseNumber(source[j:], 4, b); err != nil {
					return
				}
				j += diff
			case 'm':
				if month, diff, err = parseNumber(source[j:], 2, 'm'); err != nil {
					return
				}
				j += diff
			case 'B':
				if month, diff, err = lookup(source[j:], longMonthNames, 'B'); err != nil {
					return
				}
				j += diff
			case 'b', 'h':
				if month, diff, err = lookup(source[j:], shortMonthNames, b); err != nil {
					return
				}
				j += diff
			case 'A':
				if _, diff, err = lookup(source[j:], longWeekNames, 'A'); err != nil {
					return
				}
				j += diff
			case 'a':
				if _, diff, err = lookup(source[j:], shortWeekNames, 'a'); err != nil {
					return
				}
				j += diff
			case 'w':
				if j >= l || source[j] < '0' || '6' < source[j] {
					err = parseFormatError(b)
					return
				}
				j++
			case 'u':
				if j >= l || source[j] < '1' || '7' < source[j] {
					err = parseFormatError(b)
					return
				}
				j++
			case 'V', 'U', 'W':
				if _, diff, err = parseNumber(source[j:], 2, b); err != nil {
					return
				}
				j += diff
			case 'e':
				if j < l && source[j] == ' ' {
					j++
				}
				fallthrough
			case 'd':
				if day, diff, err = parseNumber(source[j:], 2, b); err != nil {
					return
				}
				j += diff
			case 'j':
				if yday, diff, err = parseNumber(source[j:], 3, 'j'); err != nil {
					return
				}
				j += diff
			case 'k':
				if j < l && source[j] == ' ' {
					j++
				}
				fallthrough
			case 'H':
				if hour, diff, err = parseNumber(source[j:], 2, b); err != nil {
					return
				}
				j += diff
			case 'l':
				if j < l && source[j] == ' ' {
					j++
				}
				fallthrough
			case 'I':
				if hour, diff, err = parseNumber(source[j:], 2, b); err != nil {
					return
				}
				j += diff
				if hour == 12 {
					hour = 0
				}
			case 'p', 'P':
				var ampm int
				if ampm, diff, err = lookup(source[j:], []string{"AM", "PM"}, 'p'); err != nil {
					return
				}
				j += diff
				pm = ampm == 2
			case 'M':
				if min, diff, err = parseNumber(source[j:], 2, 'M'); err != nil {
					return
				}
				j += diff
			case 'S':
				if sec, diff, err = parseNumber(source[j:], 2, 'S'); err != nil {
					return
				}
				j += diff
			case 's':
				var unix int
				if unix, diff, err = parseNumber(source[j:], 10, 's'); err != nil {
					return
				}
				t = time.Unix(int64(unix), 0).In(time.UTC)
				var mon time.Month
				year, mon, day = t.Date()
				hour, min, sec = t.Clock()
				month = int(mon)
				j += diff
			case 'f':
				var msec int
				if msec, diff, err = parseNumber(source[j:], 6, 'f'); err != nil {
					return
				}
				j += diff
				nsec = msec * 1000
				for diff < 6 {
					nsec *= 10
					diff++
				}
			case 'Z':
				k := j
				for ; k < l; k++ {
					if c := source[k]; c < 'A' || 'Z' < c {
						break
					}
				}
				t, err = time.Parse("MST", source[j:k])
				if err != nil {
					return
				}
				loc = t.Location()
				j = k
			case 'z':
				if j+5 > l {
					err = parseFormatError(b)
					return
				}
				var offset int
				switch source[j] {
				case '+', '-':
					var hour, min int
					if hour, err = strconv.Atoi(source[j+1 : j+3]); err != nil {
						return
					}
					if min, err = strconv.Atoi(source[j+3 : j+5]); err != nil {
						return
					}
					offset = (hour*60 + min) * 60
					if source[j] == '-' {
						offset = -offset
					}
				default:
					err = parseFormatError(b)
					return
				}
				loc = time.FixedZone("", offset)
				j += 5
			case 't', 'n':
				k := j
			K:
				for ; k < l; k++ {
					switch source[k] {
					case ' ', '\t', '\n', '\v', '\f', '\r':
					default:
						break K
					}
				}
				if k == j {
					err = errors.New("expected a space")
					return
				}
				j = k
			case '%':
				if j >= l || source[j] != b {
					err = fmt.Errorf("expected %q", b)
					return
				}
				j++
			default:
				if pending == "" {
					var ok bool
					if pending, ok = compositions[b]; ok {
						break
					}
					err = fmt.Errorf(`unexpected format: "%%%c"`, b)
					return
				}
				if j >= l || source[j] != b {
					err = fmt.Errorf("expected %q", b)
					return
				}
				j++
			}
			if pending != "" {
				b, pending = pending[0], pending[1:]
				goto L
			}
		} else if j >= len(source) || source[j] != b {
			err = fmt.Errorf("expected %q", b)
			return
		} else {
			j++
		}
	}
	if j < len(source) {
		err = fmt.Errorf("unconverted string: %q", source[j:])
		return
	}
	if pm {
		hour += 12
	}
	if century > 0 {
		year = century*100 + year%100
	}
	if yday > 0 {
		return time.Date(year, time.January, 1, hour, min, sec, nsec, loc).AddDate(0, 0, yday-1), nil
	}
	return time.Date(year, time.Month(month), day, hour, min, sec, nsec, loc), nil
}

func isDigit(c byte) bool {
	return '0' <= c && c <= '9'
}

type parseFormatError byte

func (err parseFormatError) Error() string {
	return fmt.Sprintf("cannot parse %%%c", byte(err))
}

func parseNumber(source string, max int, format byte) (int, int, error) {
	if len(source) > 0 && isDigit(source[0]) {
		for i := 1; i < max; i++ {
			if i >= len(source) || !isDigit(source[i]) {
				val, err := strconv.Atoi(string(source[:i]))
				return val, i, err
			}
		}
		val, err := strconv.Atoi(string(source[:max]))
		return val, max, err
	}
	return 0, 0, parseFormatError(format)
}

func lookup(source string, candidates []string, format byte) (int, int, error) {
L:
	for i, xs := range candidates {
		for j, x := range []byte(xs) {
			if j >= len(source) {
				continue L
			}
			if y := source[j]; x != y && x|('a'-'A') != y|('a'-'A') {
				continue L
			}
		}
		return i + 1, len(xs), nil
	}
	return 0, 0, parseFormatError(format)
}
