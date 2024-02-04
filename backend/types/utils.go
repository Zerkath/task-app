package types

import (
    "strconv"
)

func ToInt(s string, defaultValue int) int {
    i, err := strconv.Atoi(s)
    if err != nil {
        return defaultValue
    }
    return i
}
