package main

import (
	//"fmt"
	"encoding/csv"
	"github.com/gocolly/colly/v2"
	"strings"
	"log"
	"os"
	"net/http"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/credentials"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/joho/godotenv"
)

func main() {

	fmt.Println("Starting the scraper")
	//Create a CSV file
	fName := GetValue("FILE")
	fmt.Println(fName)
	file , err := os.Create(fName)
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()
	
	writer.Write([]string{"Instruction", "URL"})

    // List of URLs to scrape
    urls := []string{
        "https://docs.digitalocean.com/products/app-platform/",
		"https://docs.digitalocean.com/reference/doctl/reference/apps/",
    }

	c := colly.NewCollector(
		colly.AllowedDomains("docs.digitalocean.com"),
	)

	//On every a element which has href attribute call callback
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")

		if !strings.HasPrefix(link, "/products/app-platform") && !strings.HasPrefix(link, "/reference/doctl/reference/apps") {
			return
		}
		// start scraping the page under the link found
		//fmt.Println(link)
		e.Request.Visit(link)
	})

	c.OnHTML(`div[id=header-subheader]`, func(e *colly.HTMLElement) {
		log.Println("Doc found", e.Request.URL)
		resp, err := http.Get(fmt.Sprintf("%v",e.Request.URL))
		if err != nil {
			log.Fatal("Cannot get the page", err)
		}

		log.Printf("Response is %d", resp.StatusCode)

		if resp.StatusCode == 200 {
			title := strings.Split(e.ChildText("h1"), "\n")[0]
			log.Println(title)
			writer.Write([]string{title, e.Request.URL.String()})
		}
	})
	for _, url := range urls{
    	c.Visit(url)
	}
	//Upload file to s3
	UploadToS3(GetValue("FILE"))
	fmt.Println("End of Scrapper")
}


func UploadToS3(filename string) (string , error) {

	var name string = strings.TrimSuffix(filename, ".jpg") 
	fmt.Println(name)
    key := GetValue("SPACES_KEY")
    secret := GetValue("SPACES_SECRET")
	endpoint := GetValue("SPACES_ENDPOINT")
	if key == "" || secret == "" || endpoint == "" {
		log.Fatal("Missing S3 Credentials")
	}
    s3Config := &aws.Config{
        Credentials: credentials.NewStaticCredentials(key, secret, ""),
        Endpoint:    aws.String(endpoint),
        Region:      aws.String("nyc3"),
        S3ForcePathStyle: aws.Bool(false), // // Configures to use subdomain/virtual calling format. Depending on your version, alternatively use o.UsePathStyle = false
    }

    newSession := session.New(s3Config)
    s3Client := s3.New(newSession)
	fmt.Printf("%#v", s3Client)
	uploader := s3manager.NewUploader(newSession)

	pwd, _ := os.Getwd()

    filepath := pwd + "/" + filename
	fmt.Println(filepath)
	f, err := os.Open(filepath)
	defer f.Close()
	if err != nil {
		fmt.Printf("error: %v", err)
	}

	// Upload the file to S3.
	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String("backend"),
		Key:    aws.String(name),
		Body:   f,
	})
	if err != nil {	
		fmt.Printf("error while uploading photos: %v \n", err)
		return "", err
	}
	fmt.Printf("Result: %#v", result)
	fmt.Printf("file uploaded to, %s\n", aws.StringValue(&result.Location))
	
	return result.Location, nil
}

// GetValue returns configuration value based on a given key from the .env file
func GetValue(key string) string {
	fmt.Println(os.Getenv("GO_ENV"))
	env := os.Getenv("GO_ENV")
    // load the .env file
	fmt.Printf("The env value is %s \n", env)

	if os.Getenv("GO_ENV") != "PRODUCTION" {
		err := godotenv.Load(".env")
		if err != nil {
			log.Fatalf("Error loading .env file!!\n")
		}
	}

    // return the value based on a given key
	return os.Getenv(key)
}