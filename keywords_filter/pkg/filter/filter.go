package filter

import (
	"bufio"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/importcjj/sensitive"
)

type IFilter interface {
	Validate(text string) (bool, string)
	FindAll(text string) []string
}

type filter struct {
	filter *sensitive.Filter
}

func (f *filter) Validate(text string) (bool, string) {
	return f.filter.Validate(text)
}

func (f *filter) FindAll(text string) []string {
	return f.filter.FindAll(text)
}

var _filter *filter

func Getfilter() *filter {
	return _filter
}

func InitFilter(dicFilelog string) {
	if dicFilelog != "" {
		_, err := os.Stat(dicFilelog)
		if os.IsNotExist(err) {
			log.Fatalf("字典文件 %s 未找到: %v", dicFilelog, err)
		}
		f := sensitive.New()
		f.UpdateNoisePattern("")
		f.LoadWordDict(dicFilelog)
		_filter = &filter{filter: f}
		return
	}
	log.Fatal("配置文件未找到")
}

func OverWriteDict(dicFilepath string) error {
	file, err := os.OpenFile(dicFilepath, os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer file.Close()
	keymap := make(map[string]struct{})
	scanner := bufio.NewScanner(file)
	re := regexp.MustCompile(`\p{Han}+`)
	newcontent := ""
	for scanner.Scan() {
		words := scanner.Text()
		words = strings.Trim(words, " ")
		if _, ok := keymap[words]; ok {
			continue
		}
		keymap[words] = struct{}{}
		matchword := re.FindString(words)
		if matchword == "" {
			newcontent += " " + words + " \n"
		} else {
			newcontent += words + "\n"
		}
	}
	newcontent = strings.Trim(newcontent, "\n")
	_, err = file.Write([]byte(newcontent))
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}
