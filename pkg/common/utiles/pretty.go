package utiles

import "encoding/json"

func Pretty(in any) string {
	buf, err := json.MarshalIndent(in, "", "  ")
	if err != nil {
		return " \\(-_-)/ "
	}
	return string(buf)
}
