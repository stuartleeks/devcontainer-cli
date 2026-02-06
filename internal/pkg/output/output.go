package output

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/mattn/go-colorable"
	"github.com/neilotoole/jsoncolor"

	"github.com/spf13/cobra"

	jmespath "github.com/jmespath/go-jmespath"
)

type OutputFormat string

const (
	OutputJson      OutputFormat = "json"
	OutputJsonC     OutputFormat = "jsonc"
	OutputTsv       OutputFormat = "tsv"
	OutputTable     OutputFormat = "table"
	OutputTableBody OutputFormat = "tablebody"
	// TODO - does raw still make sense? offer value?
	// Raw   OutputFormat = "raw"
)

func AddOutputAndQueryFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("output", "o", "table", "output format")
	cmd.Flags().StringP("query", "q", "", "JMESPath query to apply to the result")
}
func GetOutputAndQueryValues(cmd *cobra.Command, defaultQuery string) (string, string, error) {
	outputFormat, err := cmd.Flags().GetString("output")
	if err != nil {
		return "", "", err
	}
	query, err := cmd.Flags().GetString("query")
	if err != nil {
		return "", "", err
	}
	if query == "" {
		query = defaultQuery
	}
	return outputFormat, query, err
}

// RoundtripMarshal marshals and unmarshals the data
// so that properties match the JSON property names
func RoundtripMarshal(data interface{}) ([]interface{}, error) {
	buf, err := json.Marshal(data)
	if err != nil {
		return []interface{}{}, fmt.Errorf("error marshalling result: %s", err)
	}
	var result []interface{}
	err = json.Unmarshal(buf, &result)
	if err != nil {
		return []interface{}{}, fmt.Errorf("error unmarshalling result: %s", err)
	}
	return result, nil
}

type OutputLink struct {
	Text string
	Link string
}

func (ol OutputLink) Encode() string {
	return fmt.Sprintf("||OutputLink||%s||%s", ol.Text, ol.Link)
}
func DecodeOutputLink(encodedValue string) *OutputLink {
	if !strings.HasPrefix(encodedValue, "||OutputLink||") {
		return nil
	}
	parts := strings.Split(strings.TrimPrefix(encodedValue, "||OutputLink||"), "||")
	if len(parts) != 2 {
		return nil
	}
	return &OutputLink{
		Text: parts[0],
		Link: parts[1],
	}
}

func GetTerminalLink(url string, displayText string) string {
	return fmt.Sprintf("\033]8;;%s\033\\%s\033]8;;\033\\", url, displayText)
}

// func formatEscape(s string) string {
// 	return strings.ReplaceAll(s, "%", "%%")
// }

//	func formatterFromColor(color *color.Color) func(int, int, interface{}) string {
//		colorFmt := color.SprintfFunc()
//		return func(index int, width int, val interface{}) string {
//			return colorFmt("%-*s", width, val)
//		}
//	}
func firstColumnFormatterFromColor(color *color.Color) func(int, int, interface{}) string {
	colorFmt := color.SprintfFunc()
	return func(index int, width int, val interface{}) string {
		if index == 0 {
			return colorFmt("%-*s", width, val)
		}
		return fmt.Sprintf("%-*s", width, val)
	}
}
func newLinkFormatter(inner Formatter) func(int, int, interface{}) string {
	return func(index int, width int, val interface{}) string {
		var valText string
		switch v := val.(type) {
		case int:
			valText = fmt.Sprintf("%d", v)
		case float64:
			if v == float64(int(v)) {
				valText = fmt.Sprintf("%.0f", v)
			} else {
				valText = fmt.Sprintf("%.1f", v)
			}
		case string:
			valText = v
		case OutputLink:
			amountToPad := width - len(v.Text)
			valText = GetTerminalLink(v.Link, v.Text) + strings.Repeat(" ", amountToPad)
		default:
			valText = fmt.Sprintf("%v", v)
		}

		return inner(index, width, valText)
	}
}

func newHeaderFormatter(color *color.Color) func(int, int, interface{}) string {
	// TODO take color as an input and return a formatter function
	colorFmt := color.SprintfFunc()
	return func(index int, width int, val interface{}) string {
		return colorFmt("%s", fmt.Sprintf("%-*s", width, val))
	}

}

func newLinkWidthFunc() WidthFunc {
	return func(value interface{}) int {
		if link, ok := value.(OutputLink); ok {
			return DefaultWidthFunc(link.Text)
		}
		return DefaultWidthFunc(value)
	}
}

func OutputResult(w io.Writer, result interface{}, outputFormat string, query string, defaultTableFields []string) error {

	if outputFormat == "" {
		outputFormat = string(OutputJsonC)
	}

	output := result

	tableFields := []string{}
	if query != "" {
		var err error
		output, err = jmespath.Search(query, output)
		if err != nil {
			return fmt.Errorf("error executing query: %s", err)
		}
		parser := jmespath.NewParser()
		ast, err := parser.Parse(query)
		if err != nil {
			return fmt.Errorf("error parsing query: %s", err)
		}
		tableFields, err = ast.GetResultFields()
		if err != nil {
			return fmt.Errorf("error getting query result fields: %s", err)
		}
	}
	if len(tableFields) == 0 {
		tableFields = defaultTableFields
	}

	switch outputFormat {
	case string(OutputJson):
		outputJson, _ := json.Marshal(output)
		_, _ = fmt.Fprintln(w, string(outputJson))
	case string(OutputJsonC):
		// formattedJson, _ := json.MarshalIndent(output, "", "  ")
		var enc *jsoncolor.Encoder
		if jsoncolor.IsColorTerminal(w) {
			// Safe to use color
			var out io.Writer
			if f, ok := w.(*os.File); ok {
				out = colorable.NewColorable(f) // needed for Windows
			} else {
				out = w
			}
			enc = jsoncolor.NewEncoder(out)
			// DefaultColors are similar to jq
			clrs := jsoncolor.DefaultColors()

			enc.SetColors(clrs)
		} else {
			// Can't use color; but the encoder will still work
			enc = jsoncolor.NewEncoder(w)
		}
		enc.SetIndent("", "  ")
		// TODO: Implement JSON syntax highlighting if needed
		err := enc.Encode(output)
		if err != nil {
			return fmt.Errorf("error encoding JSON: %s", err)
		}
	// case string(Raw):
	// 	fmt.Fprintln(w, output)
	case string(OutputTsv):
		// Assuming output is a slice of maps
		rows, ok := output.([]map[string]interface{})
		if ok {
			for _, row := range rows {
				values := make([]string, 0, len(row))
				for _, value := range row {
					values = append(values, fmt.Sprint(value))
				}
				_, _ = fmt.Fprintln(w, strings.Join(values, "\t"))
			}
		} else {
			_, _ = fmt.Fprintln(w, output)
		}
	case string(OutputTable), string(OutputTableBody):
		if output == nil {
			return nil
		}

		var ok bool
		var outputList []interface{}
		if outputList, ok = output.([]interface{}); !ok {
			return fmt.Errorf("output is not a list")
		}

		tmpTableFields := make([]interface{}, len(tableFields)) // convert to []interface{}
		for i, field := range tableFields {
			tmpTableFields[i] = field
		}

		tbl := NewTable(tmpTableFields...).WithWriter(w)
		valueFmt := newLinkFormatter(firstColumnFormatterFromColor(color.New(color.FgYellow)))
		widthFunc := newLinkWidthFunc()
		headerFmt := newHeaderFormatter(color.New(color.FgGreen, color.Bold))
		tbl.WithHeaderFormatter(headerFmt).WithValueFormatter(valueFmt).WithWidthFunc(widthFunc)
		for _, listItem := range outputList {
			var itemMap map[string]interface{}
			if itemMap, ok = listItem.(map[string]interface{}); !ok {
				return fmt.Errorf("output is not a list of maps")
			}
			values := make([]interface{}, len(tableFields))
			for i := 0; i < len(tableFields); i++ {
				value := itemMap[tableFields[i]]
				if tmpMap, ok := value.(map[string]interface{}); ok {
					if len(tmpMap) == 2 && tmpMap["Text"] != nil && tmpMap["Link"] != nil {
						// we have a round-tripped OutputLink
						values[i] = OutputLink{Text: tmpMap["Text"].(string), Link: tmpMap["Link"].(string)} //.Encode()
					}
				} else {
					values[i] = value
				}
			}
			tbl.AddRow(values...)
		}
		if outputFormat == string(OutputTable) {
			if err := tbl.Print(); err != nil {
				return err
			}
		} else {
			if err := tbl.PrintRows(); err != nil {
				return err
			}
		}

	default:
		return fmt.Errorf("unhandled output format: '%s'", outputFormat)
	}
	return nil
}
