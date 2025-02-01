//go:build !wasm
// +build !wasm

package term

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/refaktor/rye/env"
	"github.com/refaktor/rye/util"

	"github.com/pkg/term"
)

func DisplayBlock(bloc env.Block, idx *env.Idxs) (env.Object, bool) {
	HideCur()
	curr := 0
	moveUp := 0
	mode := 0 // 0 - human, 1 - dev
	len := bloc.Series.Len()
DODO:
	if moveUp > 0 {
		CurUp(moveUp)
	}
	SaveCurPos()
	for i, v := range bloc.Series.S {
		ClearLine()
		if i == curr {
			ColorBrGreen()
			Bold()
			fmt.Print("\u00bb ")
		} else {
			fmt.Print(" ")
		}
		switch ob := v.(type) {
		case env.Object:
			if mode == 0 {
				fmt.Println("" + ob.Print(*idx) + "")
			} else {
				fmt.Println("" + ob.Inspect(*idx) + "")
			}
		default:
			fmt.Println("" + fmt.Sprint(ob) + "")
		}
		CloseProps()
		// term.CurUp(1)
	}

	moveUp = bloc.Series.Len()

	defer func() {
		// Show cursor.
		fmt.Printf("\033[?25h")
	}()

	// RestoreCurPos()

	for {
		ascii, keyCode, err := GetChar()

		if (ascii == 3 || ascii == 27) || err != nil {
			//fmt.Println()
			ShowCur()
			return nil, true
		}

		if ascii == 13 {
			//fmt.Println()
			return bloc.Series.Get(curr), false
		}

		if ascii == 77 || ascii == 109 {
			if mode == 0 {
				mode = 1
			} else {
				mode = 0
			}
			goto DODO
		}

		if keyCode == 40 {
			curr++
			if curr > len-1 {
				curr = 0
			}
			goto DODO
		} else if keyCode == 38 {
			curr--
			if curr < 0 {
				curr = len - 1
			}
			goto DODO
		}
	}
}

func DisplaySelection(bloc env.Block, idx *env.Idxs, right int) (env.Object, bool) {
	HideCur()
	curr := 0
	moveUp := 0
	mode := 0 // 0 - human, 1 - dev
	len := bloc.Series.Len() / 2
DODO1:
	if moveUp > 0 {
		CurUp(moveUp)
	}
	SaveCurPos()
	idents := make([]int, (bloc.Series.Len()/2)+3)
	//fmt.Println("---")
	//fmt.Println(bloc.Series.Len())
	for i := 0; i < bloc.Series.Len(); i += 2 {
		//fmt.Println(i)
		// todo check if it's word
		idents[i/2] = bloc.Series.Get(i).(env.Word).Index
		label := bloc.Series.Get(i + 1)
		// ClearLine()
		CurRight(right)
		if i/2 == curr {
			ColorMagenta()
			Bold()
			fmt.Print("\u00bb ")
		} else {
			fmt.Print(" ")
		}
		switch ob := label.(type) {
		case env.String:
			fmt.Println(ob.Value + " ")
		default:
			fmt.Println("" + fmt.Sprint(ob) + "***")
		}
		CloseProps()
		// term.CurUp(1)
	}

	moveUp = bloc.Series.Len() / 2

	defer func() {
		// Show cursor.
		fmt.Printf("\033[?25h")
	}()

	// RestoreCurPos()

	for {
		ascii, keyCode, err := GetChar()

		if (ascii == 3 || ascii == 27) || err != nil {
			//fmt.Println()
			ShowCur()
			return nil, true
		}

		if ascii == 13 {
			fmt.Println()
			return *env.NewWord(idents[curr]), false
		}

		if ascii == 77 || ascii == 109 {
			if mode == 0 {
				mode = 1
			} else {
				mode = 0
			}
			goto DODO1
		}

		if keyCode == 40 {
			curr++
			if curr > len-1 {
				curr = 0
			}
			goto DODO1
		} else if keyCode == 38 {
			curr--
			if curr < 0 {
				curr = len - 1
			}
			goto DODO1
		}
	}
}

func DisplayInputField(right int, mlen int) (env.Object, bool) {
	// HideCur()
	//curr := 0
	moveUp := 0
	text := ""
	//DODO1:
	if moveUp > 0 {
		CurUp(moveUp)
	}

	defer func() {
		// Show cursor.
		// fmt.Printf("\033[?25h")
	}()

	// RestoreCurPos()
	//fmt.Println(".")
	//fmt.Println(".")
	//CurUp(2)

	CurRight(right)

	SaveCurPos()

	for {
		letter, ascii, _, err := GetChar2()
		//		letter := fmt.Scan()
		// RestoreCurPos()
		//CurDown(1)
		//fmt.Print("-----------")
		//fmt.Print(ascii)
		//CurUp(1)
		if (ascii == 3 || ascii == 27) || err != nil {
			// ShowCur()
			return nil, true
		} else if (ascii == 127) || err != nil {
			text = text[0 : len(text)-1]
			RestoreCurPos()
			fmt.Print("                  ")
			RestoreCurPos()
			fmt.Print(text)
		} else if ascii == 13 {
			fmt.Println("")
			fmt.Println("")
			return *env.NewString(text), false
		} else {
			if len(text) < mlen {
				text += letter
				RestoreCurPos()
				fmt.Print(text)
			}
		}

		// if keyCode == 40 {
		// }
		// else if keyCode == 38 {
		// }
	}
}

func DisplayDict(bloc env.Dict, idx *env.Idxs) (env.Object, bool) {
	HideCur()
	curr := 0
	moveUp := 0
	mode := 0 // 0 - human, 1 - dev
	len1 := len(bloc.Data)
	// make a slice for keys
	keys := make([]string, len(bloc.Data))
	i := 0
	for k := range bloc.Data {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
DODO:
	if moveUp > 0 {
		CurUp(moveUp)
	}
	SaveCurPos()
	for ii, k := range keys {
		// for k, v := range bloc.Data {
		v := bloc.Data[k]
		ClearLine()
		if ii == curr {
			ColorBrGreen()
			Bold()
			fmt.Print("\u00bb ")
		} else {
			fmt.Print(" ")
		}
		Bold()
		fmt.Print(k + ": ")
		ResetBold()
		switch ob := v.(type) {
		case env.Object:
			if mode == 0 {
				fmt.Println("" + ob.Print(*idx) + "")
			} else {
				fmt.Println("" + ob.Inspect(*idx) + "")
			}
		default:
			fmt.Println(" " + fmt.Sprint(ob) + " ")
		}
		CloseProps()
		// term.CurUp(1)
	}

	moveUp = len(bloc.Data)

	defer func() {
		// Show cursor.
		fmt.Printf("\033[?25h")
	}()

	// RestoreCurPos()

	for {
		ascii, keyCode, err := GetChar()

		if (ascii == 3 || ascii == 27) || err != nil {
			//fmt.Println()
			ShowCur()
			return nil, true
		}

		if ascii == 13 {
			//fmt.Println()
			ret := ""
			for ii, k := range keys {
				if ii == curr {
					ret = k
				}
			}
			return *env.NewString(ret), false // bloc.Series.Get(curr), false
		}

		if ascii == 120 {
			//fmt.Println()
			var ret env.Object
			for ii, k := range keys {
				if ii == curr {
					ret = bloc.Data[k].(env.Object)
				}
			}
			return ret, false // bloc.Series.Get(curr), false
		}

		if ascii == 77 || ascii == 109 {
			if mode == 0 {
				mode = 1
			} else {
				mode = 0
			}
			goto DODO
		}

		if keyCode == 40 {
			curr++
			if curr > len1-1 {
				curr = 0
			}
			goto DODO
		} else if keyCode == 38 {
			curr--
			if curr < 0 {
				curr = len1 - 1
			}
			goto DODO
		}
	}
}

func DisplayTableRow(bloc env.TableRow, idx *env.Idxs) (env.Object, bool) {
	HideCur()
	curr := 0
	moveUp := 0
	mode := 0 // 0 - human, 1 - dev
	len1 := len(bloc.Values)
	// make a slice for keys
	/* keys := make([]string, len(bloc.Data))
	i := 0
	for k, _ := range bloc.Data {
		keys[i] = k
		i++
	}
	sort.Strings(keys)*/
DODO:
	if moveUp > 0 {
		CurUp(moveUp)
	}
	SaveCurPos()
	for ii, k := range bloc.Uplink.Cols {
		// for k, v := range bloc.Data {
		v := bloc.Values[ii]
		ClearLine()
		if ii == curr {
			ColorBrGreen()
			Bold()
			fmt.Print("\u00bb ")
		} else {
			fmt.Print(" ")
		}
		Bold()
		fmt.Print(k + ": ")
		ResetBold()
		switch ob := v.(type) {
		case env.Object:
			if mode == 0 {
				fmt.Println("" + ob.Print(*idx) + "")
			} else {
				fmt.Println("" + ob.Inspect(*idx) + "")
			}
		default:
			fmt.Println(" " + fmt.Sprint(ob) + " ")
		}
		CloseProps()
		// term.CurUp(1)
	}

	moveUp = len(bloc.Values)

	defer func() {
		// Show cursor.
		fmt.Printf("\033[?25h")
	}()

	// RestoreCurPos()

	for {
		ascii, keyCode, err := GetChar()

		if (ascii == 3 || ascii == 27) || err != nil {
			//fmt.Println()
			ShowCur()
			return nil, true
		}

		if ascii == 13 {
			return env.ToRyeValue(bloc.Values[curr]), false
		}

		if ascii == 120 {
			//fmt.Println()
			return env.String{Value: bloc.Uplink.Cols[curr]}, false
		}

		if ascii == 77 || ascii == 109 {
			if mode == 0 {
				mode = 1
			} else {
				mode = 0
			}
			goto DODO
		}

		if keyCode == 40 {
			curr++
			if curr > len1-1 {
				curr = 0
			}
			goto DODO
		} else if keyCode == 38 {
			curr--
			if curr < 0 {
				curr = len1 - 1
			}
			goto DODO
		}
	}
}

func DisplayTable(bloc env.Table, idx *env.Idxs) (env.Object, bool) {
	HideCur()
	curr := 0
	moveUp := 0
	mode := 0 // 0 - human, 1 - dev
	// get the ideal widths of columns
	widths := make([]int, len(bloc.Cols))
	// check all col names
	for ic, col := range bloc.Cols {
		widths[ic] = len(col) + 1
	}
	// check all data
	for _, r := range bloc.Rows {
		for ic, v := range r.Values {
			ww := 5
			switch val := v.(type) {
			case string:
				ww = len(val) + 2
				if ww > 52 {
					ww = 52
				}
			case int64:
				ww = len(strconv.Itoa(int(val))) + 1
			case env.Integer:
				ww = len(strconv.Itoa(int(val.Value))) + 1
			case float64:
				ww = len(strconv.FormatFloat(val, 'f', 2, 64)) + 1
			case env.Decimal:
				ww = len(strconv.FormatFloat(val.Value, 'f', 2, 64)) + 1
			case env.String:
				ww = len(val.Print(*idx))
				if ww > 52 {
					ww = 52
				}
				//if ww > 60 {
				// ww = 60
				//}
			case env.Vector:
				ww = len(val.Print(*idx))
			}
			if len(widths) > ic && widths[ic] < ww {
				widths[ic] = ww + 1
			}
		}
	}
	fulwidth := 0
	for _, w := range widths {
		fulwidth += w + 2
	}

DODO:
	if moveUp > 0 {
		CurUp(moveUp)
	}
	SaveCurPos()
	for ic, cn := range bloc.Cols {
		Bold()
		fmt.Printf("| %-"+strconv.Itoa(widths[ic])+"s", cn)
		CloseProps()
	}
	fmt.Println("|")
	fmt.Println("+" + strings.Repeat("-", fulwidth-1) + "+")

	for range bloc.Rows {
		ClearLine()
	}
	for i, r := range bloc.Rows {
		if i == curr {
			ColorBrGreen()
			fmt.Print("")
		} else {
			fmt.Print("")
		}
		for ic, v := range r.Values {
			if ic < len(widths) {
				// fmt.Println(v)
				switch ob := v.(type) {
				case env.Object:
					if mode == 0 {
						fmt.Printf("| %-"+strconv.Itoa(widths[ic])+"s", util.TruncateString(ob.Print(*idx), widths[ic]))
						//fmt.Print("| " +  + "\t")
					} else {
						fmt.Printf("| %-"+strconv.Itoa(widths[ic])+"s", ob.Inspect(*idx))
						//fmt.Print("| " +  + "\t")
					}
				default:
					fmt.Printf("| %-"+strconv.Itoa(widths[ic])+"s", fmt.Sprint(ob))
					///fmt.Print("| " + +"\t")
				}
			}
			// term.CurUp(1)
		}
		CloseProps()
		fmt.Println("|")
	}

	moveUp = len(bloc.Rows) + 2

	defer func() {
		// Show cursor.
		fmt.Printf("\033[?25h")
	}()

	// RestoreCurPos()

	for {
		ascii, keyCode, err := GetChar()

		if (ascii == 3 || ascii == 27) || err != nil {
			fmt.Println()
			ShowCur()
			return nil, true
		}

		if ascii == 13 {
			fmt.Println()
			return bloc.GetRowNew(curr), false // bloc.Series.Get(curr), false
		}

		if ascii == 77 || ascii == 109 {
			if mode == 0 {
				mode = 1
			} else {
				mode = 0
			}
			goto DODO
		}

		if keyCode == 40 {
			curr++
			goto DODO
		} else if keyCode == 38 {
			curr--
			goto DODO
		}
	}
}

// ideation:
// .display\custom fn { x } { -> 'subject .elipsis 20 .red .prn , spacer 2 , -> 'score .align-right 10 .print }
func DisplayTableCustom(bloc env.Table, myfn func(row env.Object, iscurr env.Integer), idx *env.Idxs) (env.Object, bool) {
	HideCur()
	curr := 0
	moveUp := 0
	mode := 0 // 0 - human, 1 - dev
	// get the ideal widths of columns
	widths := make([]int, len(bloc.Cols))
	// check all col names
	for ic, col := range bloc.Cols {
		widths[ic] = len(col) + 1
	}
	// check all data
	for _, r := range bloc.Rows {
		for ic, v := range r.Values {
			ww := 5
			switch val := v.(type) {
			case string:
				ww = len(val) + 2
				if ww > 52 {
					ww = 52
				}
			case int64:
				ww = len(strconv.Itoa(int(val))) + 1
			case env.Integer:
				ww = len(strconv.Itoa(int(val.Value))) + 1
			case float64:
				ww = len(strconv.FormatFloat(val, 'f', 2, 64)) + 1
			case env.Decimal:
				ww = len(strconv.FormatFloat(val.Value, 'f', 2, 64)) + 1
			case env.String:
				ww = len(val.Print(*idx))
				if ww > 52 {
					ww = 52
				}
				//if ww > 60 {
				// ww = 60
				//}
			case env.Vector:
				ww = len(val.Print(*idx))
			}
			if len(widths) > ic && widths[ic] < ww {
				widths[ic] = ww + 1
			}
		}
	}
	fulwidth := 0
	for _, w := range widths {
		fulwidth += w + 2
	}

DODO:
	if moveUp > 0 {
		CurUp(moveUp)
	}
	SaveCurPos()
	/* for ic, cn := range bloc.Cols {
		Bold()
		fmt.Printf("| %-"+strconv.Itoa(widths[ic])+"s", cn)
		CloseProps()
	}
	fmt.Println("|")
	fmt.Println("+" + strings.Repeat("-", fulwidth-1) + "+")
	*/

	for range bloc.Rows {
		ClearLine()
	}
	for i, r := range bloc.Rows {
		iscurr := *env.NewInteger(0)
		if i == curr {
			iscurr = *env.NewInteger(1)
		}

		// call funtion with row and is-current value
		myfn(r, iscurr)

		//CloseProps()
		// fmt.Println("|")
	}

	moveUp = len(bloc.Rows)

	defer func() {
		// Show cursor.
		fmt.Printf("\033[?25h")
	}()

	// RestoreCurPos()

	for {
		ascii, keyCode, err := GetChar()

		if (ascii == 3 || ascii == 27) || err != nil {
			fmt.Println()
			ShowCur()
			return nil, true
		}

		if ascii == 13 {
			fmt.Println()
			return bloc.GetRowNew(curr), false // bloc.Series.Get(curr), false
		}

		if ascii == 77 || ascii == 109 {
			if mode == 0 {
				mode = 1
			} else {
				mode = 0
			}
			goto DODO
		}

		if keyCode == 40 {
			curr++
			goto DODO
		} else if keyCode == 38 {
			curr--
			goto DODO
		}
	}
}

func itoa(i int) {
	panic("unimplemented")
}

func ShowCur() {
	fmt.Print("\x1b[?25h")
}

func HideCur() {
	fmt.Print("\x1b[?25l")
}
func SaveCurPos() {
	fmt.Print("\x1b7")
}

func ClearLine() {
	fmt.Print("\x1b[0K")
}

func RestoreCurPos() {
	fmt.Print("\x1b8")
}

// Standard colors
func ColorBlack() {
	fmt.Printf("\x1b[30m")
}
func ColorRed() {
	fmt.Printf("\x1b[31m")
}
func ColorGreen() {
	fmt.Printf("\x1b[32m")
}
func ColorYellow() {
	fmt.Printf("\x1b[33m")
}
func ColorBlue() {
	fmt.Printf("\x1b[34m")
}
func ColorMagenta() {
	fmt.Printf("\x1b[35m")
}
func ColorCyan() {
	fmt.Printf("\x1b[36m")
}
func ColorWhite() {
	fmt.Printf("\x1b[37m")
}

// Standard colors returned
func StrColorBlack() string {
	return "\x1b[30m"
}
func StrColorRed() string {
	return "\x1b[31m"
}
func StrColorGreen() string {
	return "\x1b[32m"
}
func StrColorYellow() string {
	return "\x1b[33m"
}
func StrColorBlue() string {
	return "\x1b[34m"
}
func StrColorMagenta() string {
	return "\x1b[35m"
}
func StrColorCyan() string {
	return "\x1b[36m"
}
func StrColorWhite() string {
	return "\x1b[37m"
}

func StrColorBrBlack() string {
	return "\x1b[30;1m"
}

// Bright colors
func ColorBrBlack() {
	fmt.Printf("\x1b[30;1m")
}
func ColorBrRed() {
	fmt.Printf("\x1b[31;1m")
}
func ColorBrGreen() {
	fmt.Printf("\x1b[32;1m")
}
func ColorBrYellow() {
	fmt.Printf("\x1b[33;1m")
}
func ColorBrBlue() {
	fmt.Printf("\x1b[34;1m")
}
func ColorBrMagenta() {
	fmt.Printf("\x1b[36;1m")
}
func ColorBrCyan() {
	fmt.Printf("\x1b[37;1m")
}
func ColorBrWhite() {
	fmt.Printf("\x1b[37;1m")
}

// Background
func ColorBgBlack() {
	fmt.Printf("\x1b[40m")
}
func ColorBgRed() {
	fmt.Printf("\x1b[41m")
}
func ColorBgGreen() {
	fmt.Printf("\x1b[42m")
}
func ColorBgYellow() {
	fmt.Printf("\x1b[43m")
}
func ColorBgBlue() {
	fmt.Printf("\x1b[44m")
}
func ColorBgMagenta() {
	fmt.Printf("\x1b[45m")
}
func ColorBgCyan() {
	fmt.Printf("\x1b[46m")
}
func ColorBgWhite() {
	fmt.Printf("\x1b[47m")
}

// Font style
func Bold() {
	fmt.Printf("\x1b[1m")
}
func Italic() {
	fmt.Printf("\x1b[3m")
}
func Underline() {
	fmt.Printf("\x1b[4m")
}
func ResetBold() {
	fmt.Printf("\x1b[22m")
}
func CloseProps() {
	fmt.Printf("\x1b[0m")
}
func StrCloseProps() string {
	return "\x1b[0m"
}
func CurUp(n int) {
	fmt.Printf("\x1b[%dA", n)
}
func CurDown(n int) {
	fmt.Printf("\x1b[%dB", n)
}
func CurRight(n int) {
	fmt.Printf("\x1b[%dC", n)
}
func CurLeft(n int) {
	fmt.Printf("\x1b[%dD", n)
}

// From https://github.com/paulrademacher/climenu/blob/master/getchar.go

func GetChar() (ascii int, keyCode int, err error) {
	t, _ := term.Open("/dev/tty")
	if err = term.RawMode(t); err != nil {
		return
	}
	bytes := make([]byte, 3)

	var numRead int
	numRead, err = t.Read(bytes)
	if err != nil {
		return
	}
	if numRead == 3 && bytes[0] == 27 && bytes[1] == 91 {
		// Three-character control sequence, beginning with "ESC-[".

		// Since there are no ASCII codes for arrow keys, we use
		// Javascript key codes.
		if bytes[2] == 65 {
			// Up
			keyCode = 38
		} else if bytes[2] == 66 {
			// Down
			keyCode = 40
		} else if bytes[2] == 67 {
			// Right
			keyCode = 39
		} else if bytes[2] == 68 {
			// Left
			keyCode = 37
		}
	} else if numRead == 1 {
		ascii = int(bytes[0])
	}
	// else {
	// Two characters read??
	// }
	if err = t.Restore(); err != nil {
		return
	}
	t.Close()
	return
}

func GetChar2() (letter string, ascii int, keyCode int, err error) {
	t, _ := term.Open("/dev/tty")
	if err = term.RawMode(t); err != nil {
		return
	}
	bytes := make([]byte, 3)

	var numRead int
	numRead, err = t.Read(bytes)
	if err != nil {
		return
	}
	if numRead == 3 && bytes[0] == 27 && bytes[1] == 91 {
		// Three-character control sequence, beginning with "ESC-[".

		// Since there are no ASCII codes for arrow keys, we use
		// Javascript key codes.
		if bytes[2] == 65 {
			// Up
			keyCode = 38
		} else if bytes[2] == 66 {
			// Down
			keyCode = 40
		} else if bytes[2] == 67 {
			// Right
			keyCode = 39
		} else if bytes[2] == 68 {
			// Left
			keyCode = 37
		}
	} else if numRead == 1 {
		ascii = int(bytes[0])
		letter = string(bytes[0])
	} else if numRead == 2 {
		letter = string(bytes[0:2])
	} else if numRead == 3 {
		letter = string(bytes)
	}
	if err = t.Restore(); err != nil {
		return
	}
	t.Close()
	return
}
