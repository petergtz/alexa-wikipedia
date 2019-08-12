package main

import (
	"fmt"
	"os"
	"sort"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

var files = make(map[string]*os.File)

func main() {
	db := dynamodb.New(session.Must(session.NewSession(&aws.Config{
		Region:      aws.String("eu-central-1"),
		Credentials: credentials.NewStaticCredentials(os.Getenv("ACCESS_KEY_ID"), os.Getenv("SECRET_ACCESS_KEY"), ""),
	})))
	counts := make(map[string]map[string]int)
	e := db.ScanPages(&dynamodb.ScanInput{
		TableName: aws.String("AlexaWikipediaRequests"),
	}, func(output *dynamodb.ScanOutput, last bool) bool {
		for _, item := range output.Items {
			if item["RequestType"] != nil &&
				item["RequestType"].S != nil &&
				*item["RequestType"].S == "IntentRequest" &&
				item["Attributes"] != nil &&
				item["Attributes"].M != nil &&
				item["Attributes"].M["Intent"] != nil &&
				item["Attributes"].M["Intent"].S != nil &&
				*item["Attributes"].M["Intent"].S == "DefineIntent" {
				fmt.Fprintf(getFile(*item["Locale"].S), "%v\t%v\t%v\t%v\t%v\t%v\t%v\n",
					*item["UnixTimestamp"].N,
					*item["Timestamp"].S,
					*item["RequestID"].S,
					*item["UserID"].S,
					*item["SessionID"].S,
					*item["Attributes"].M["SearchQuery"].S,
					*item["Attributes"].M["ActualTitle"].S,
				)
				if counts[*item["Locale"].S] == nil {
					counts[*item["Locale"].S] = make(map[string]int)
				}
				counts[*item["Locale"].S][*item["Attributes"].M["ActualTitle"].S]++
			}
		}
		return true
	})
	PanicOnError(e)
	for locale, countsPerLocale := range counts {
		bla := make([]TitleAndCount, 0)
		for title, count := range countsPerLocale {
			bla = append(bla, TitleAndCount{title: title, count: count})
		}
		sort.Slice(bla, func(i int, j int) bool { return bla[i].count > bla[j].count })
		fmt.Println(locale, ":")
		for i, item := range bla {
			if i > 20 {
				break
			}
			fmt.Println(item.title, item.count)
		}
	}
}

type TitleAndCount struct {
	title   string
	count   int
	queries []string
}

func getFile(locale string) *os.File {
	if _, exist := files[locale]; !exist {
		var e error
		files[locale], e = os.Create("requests-" + locale)
		PanicOnError(e)
	}
	return files[locale]
}

func PanicOnError(e error) {
	if e != nil {
		panic(e)
	}
}

var Must = PanicOnError
