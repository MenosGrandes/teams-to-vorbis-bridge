package impl

import (
	"fmt"
	"sort"

	"vgc/main/engine"

	lua "github.com/yuin/gopher-lua"
)

var _ engine.Engine = (*LuaEngine)(nil)

type threshold struct {
	minScore int
	grade    float64
}

type LuaEngine struct {
	state      *lua.LState
	formula    string
	thresholds []threshold
}

func openSafeLibs(L *lua.LState) {
	for _, lib := range []struct {
		name string
		fn   lua.LGFunction
	}{
		{lua.BaseLibName, lua.OpenBase},
		{lua.TabLibName, lua.OpenTable},
		{lua.StringLibName, lua.OpenString},
		{lua.MathLibName, lua.OpenMath},
	} {
		L.Push(L.NewFunction(lib.fn))
		L.Push(lua.LString(lib.name))
		L.Call(1, 0)
	}
}

func NewLuaEngine(formula string, grades map[int]float64) (*LuaEngine, error) {
	L := lua.NewState()
	openSafeLibs(L)

	if _, err := L.LoadString(formula); err != nil {
		L.Close()
		return nil, fmt.Errorf("invalid formula: %w", err)
	}
	L.SetTop(0)

	thresholds := make([]threshold, 0, len(grades))
	for score, grade := range grades {
		thresholds = append(thresholds, threshold{minScore: score, grade: grade})
	}
	sort.Slice(thresholds, func(i, j int) bool {
		return thresholds[i].minScore > thresholds[j].minScore
	})

	return &LuaEngine{state: L, formula: formula, thresholds: thresholds}, nil
}

func (e *LuaEngine) Evaluate(cols map[string]float64) (float32, error) {
	L := e.state
	L.SetTop(0)

	for letter, val := range cols {
		L.SetGlobal(letter, lua.LNumber(val))
	}

	L.SetGlobal("getGrade", L.NewFunction(func(ls *lua.LState) int {
		score := int(ls.CheckNumber(1))
		for _, t := range e.thresholds {
			if score >= t.minScore {
				ls.Push(lua.LNumber(t.grade))
				return 1
			}
		}
		ls.Push(lua.LNumber(0))
		return 1
	}))

	fn, err := L.LoadString(e.formula)
	if err != nil {
		return 0, fmt.Errorf("loading formula: %w", err)
	}
	L.Push(fn)
	if err := L.PCall(0, 1, nil); err != nil {
		return 0, fmt.Errorf("executing formula: %w", err)
	}

	return float32(lua.LVAsNumber(L.Get(-1))), nil
}

func (e *LuaEngine) Close() {
	e.state.Close()
}
