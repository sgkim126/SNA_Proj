package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"strings"
)

func shouldNot(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	var orgFileName string
	var langDirName string
	var outputFileName string
	flag.StringVar(&orgFileName, "org", "org", "org file")
	flag.StringVar(&langDirName, "lang", "lang", "lang dir")
	flag.StringVar(&outputFileName, "output", "output", "output file name")

	flag.Parse()
	fmt.Println("Start with org:", orgFileName)
	fmt.Println("          lang:", langDirName)
	fmt.Println("        output:", outputFileName)

	outputFile, err := os.Create(outputFileName)
	shouldNot(err)
	defer outputFile.Close()
	writer := bufio.NewWriter(outputFile)

	orgs := parseOrg(orgFileName)
	langNames, langs := parseLangs(langDirName)
	lengthOfLangNames := len(langNames)

	for orgId, orgMember := range orgs {
		for i1, langName1 := range langNames {
			for i2 := i1 + 1; i2 < lengthOfLangNames; i2 += 1 {
				members1, ok := langs[langName1]
				if !ok {
					panic(errors.New(fmt.Sprintf("%s is not a valid name", langName1)))
				}

				langName2 := langNames[i2]
				members2, ok := langs[langName2]
				if !ok {
					panic(errors.New(fmt.Sprintf("%s is not a valid name", langName2)))
				}

				lang1 := filter(members1, orgMember)
				if len(lang1) == 0 {
					continue
				}
				lang2 := filter(members2, orgMember)
				if len(lang2) == 0 {
					continue
				}

				common := common_distance(lang1, lang2)
				if common == 0 {
					continue
				}
				if len(lang1)+len(lang2)-common < 5 {
					continue
				}

				writer.WriteString(fmt.Sprintf("%d %s %s", orgId, strings.Replace(langName1, " ", "-", -1), strings.Replace(langName2, " ", "-", -1)))
				writer.WriteString(fmt.Sprintf(" fraction %f", float64(len(lang1)+len(lang2)-common)/float64(len(orgMember))))
				writer.WriteString(fmt.Sprintf(" total %d", len(orgMember)))
				writer.WriteString(fmt.Sprintf(" union %d", len(lang1)+len(lang2)-common))
				writer.WriteString(fmt.Sprintf(" common %d", common))

				jaccard := jaccard_distance(lang1, lang2)
				writer.WriteString(fmt.Sprintf(" jaccard %f", jaccard))
				writer.WriteByte('\n')
				writer.Flush()
			}
		}
	}

}

func filter(lang map[int]struct{}, org []int) map[int]struct{} {
	result := make(map[int]struct{}, 0)
	for _, o := range org {
		_, ok := lang[o]
		if ok {
			result[o] = struct{}{}
		}
	}
	return result
}

func parseOrg(fileName string) map[int]([]int) {
	file, err := os.Open(fileName)
	shouldNot(err)
	defer file.Close()
	result := make(map[int]([]int), 0)
	reader := bufio.NewReader(file)
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		shouldNot(scanner.Err())

		line := scanner.Text()
		tokens := strings.Split(line, " ")

		org, err := strconv.Atoi(tokens[0])
		shouldNot(err)

		members := []int{}
		if len(tokens) < 5+2 {
			continue
		}
		for i := 2; i < len(tokens); i += 1 {
			member, err := strconv.Atoi(tokens[i])
			shouldNot(err)
			members = append(members, member)
		}
		result[org] = members
	}
	return result
}

func langName(fileName string) string {
	return fileName[0 : len(fileName)-len(".txt")]
}

func common_distance(lang1 map[int]struct{}, lang2 map[int]struct{}) int {
	c := 0
	for l1 := range lang1 {
		_, ok := lang2[l1]
		if ok {
			c += 1
		}
	}
	return c
}

func jaccard_distance(lang1 map[int]struct{}, lang2 map[int]struct{}) float64 {
	a := len(lang1)
	b := len(lang2)
	c := common_distance(lang1, lang2)
	return float64(a+b-c-c) / float64(a+b-c)
}

func parseLangs(inputDir string) ([]string, map[string](map[int]struct{})) {
	files, err := ioutil.ReadDir(inputDir)
	shouldNot(err)
	lengthOfFiles := len(files)
	members := make(map[string](map[int]struct{}), lengthOfFiles)
	langNames := make([]string, 0)
	for _, file := range files {
		name, member := parseLang(inputDir, file.Name())
		members[name] = member
		langNames = append(langNames, name)
	}
	return langNames, members
}

func parseLang(inputDir string, fileName string) (string, map[int]struct{}) {
	path := path.Join(inputDir, fileName)
	name := langName(fileName)

	input, err := os.Open(path)
	shouldNot(err)
	defer input.Close()

	reader := bufio.NewReader(input)
	scanner := bufio.NewScanner(reader)
	result := make(map[int]struct{}, 0)
	for scanner.Scan() {
		shouldNot(scanner.Err())
		id, err := strconv.Atoi(scanner.Text())
		shouldNot(err)
		result[id] = struct{}{}
	}

	return name, result
}
