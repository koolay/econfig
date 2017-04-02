package dotfile

import (
	"bufio"
	"errors"
	"fmt"
	"regexp"
	// "github.com/koolay/econfig/config"
	"os"
	"strings"
)

// ConfigItem configItem
type ConfigItem struct {
	Key     string
	Value   string
	Comment string
}

func GenerateEnvFile(configItemMap map[string]*ConfigItem, saveTo string) ([]byte, error) {
	// linesMap := make(map[string]*ConfigItem)
	return nil, nil
}

func GetConfigItemFromEnvFile(filepath string, key string) (*ConfigItem, error) {
	if items, err := ReadEnvFile(filepath); err == nil {
		if item, ok := items[key]; ok {
			return item, nil
		}
		return nil, nil
	} else {
		return nil, err
	}
}

// ReadEnvFile read .env file
func ReadEnvFile(filepath string) (map[string]*ConfigItem, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	linesMap := make(map[string]*ConfigItem)

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err = scanner.Err(); err != nil {
		return nil, err
	}

	var latestComment string
	for _, fullLine := range lines {
		if !isIgnoredLine(fullLine) {
			var configItem ConfigItem
			if isCommentLine(fullLine) {
				reg := regexp.MustCompile(`#+\s*`)
				latestComment = reg.ReplaceAllString(fullLine, "")
			} else {
				err = parseLine(&configItem, fullLine)
				if err != nil {
					return nil, err
				}
				configItem.Comment = latestComment
				latestComment = ""
				linesMap[configItem.Key] = &configItem
			}
		}
	}
	return linesMap, nil

}
func parseLine(item *ConfigItem, line string) (err error) {
	if len(line) == 0 {
		err = errors.New("zero length string")
		return
	}

	// ditch the comments (but keep quoted hashes)
	if strings.Contains(line, "#") {
		segmentsBetweenHashes := strings.Split(line, "#")
		quotesAreOpen := false
		var segmentsToKeep []string
		for _, segment := range segmentsBetweenHashes {
			if strings.Count(segment, "\"") == 1 || strings.Count(segment, "'") == 1 {
				if quotesAreOpen {
					quotesAreOpen = false
					segmentsToKeep = append(segmentsToKeep, segment)
				} else {
					quotesAreOpen = true
				}
			}

			if len(segmentsToKeep) == 0 || quotesAreOpen {
				segmentsToKeep = append(segmentsToKeep, segment)
			}
		}

		line = strings.Join(segmentsToKeep, "#")
	}

	// now split key from value
	splitString := strings.SplitN(line, "=", 2)

	if len(splitString) != 2 {
		err = errors.New("Can't separate key from value")
		return
	}

	// Parse the key
	key := splitString[0]
	if strings.HasPrefix(key, "export") {
		key = strings.TrimPrefix(key, "export")
	}
	key = strings.Trim(key, " ")

	// Parse the value
	value := splitString[1]
	// trim
	value = strings.Trim(value, " ")

	// check if we've got quoted values
	if strings.Count(value, "\"") == 2 || strings.Count(value, "'") == 2 {
		// pull the quotes off the edges
		value = strings.Trim(value, "\"'")

		// expand quotes
		value = strings.Replace(value, "\\\"", "\"", -1)
		// expand newlines
		value = strings.Replace(value, "\\n", "\n", -1)
	}

	item.Key = key
	item.Value = value

	return
}

// isCommentLine if comment line
func isCommentLine(line string) bool {
	trimmedLine := strings.Trim(line, " \n\t")
	return strings.HasPrefix(trimmedLine, "#")
}

// isIgnoredLine if empty line
func isIgnoredLine(line string) bool {
	trimmedLine := strings.Trim(line, " \n\t")
	return len(trimmedLine) == 0
}

func WriteLines(lines []string, path string, uid int, gid int) error {
	file, err := os.Create(path)

	if err != nil {
		return err
	}

	if err := os.Chown(path, uid, gid); err != nil {
		return err
	}

	defer file.Close()

	w := bufio.NewWriter(file)
	for _, line := range lines {
		fmt.Fprintln(w, line)
	}
	return w.Flush()
}
