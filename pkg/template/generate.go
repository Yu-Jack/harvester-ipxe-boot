package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
)

var (
	version     = "v1.3.0"
	outputDir   = "./public/harvester"
	templateDir = "./pkg/template"
	baseURL     = "http://192.168.1.122:3333" // ISO files are served from this URL
	primaryIP   = "192.168.122.61"
	secondaryIP = "192.168.122.62"
	token       = "123"
	config      = Config{BaseURL: baseURL, Version: version, PrimaryIP: primaryIP, SecondaryIP: secondaryIP, Token: token}
)

type Config struct {
	BaseURL     string
	Version     string
	PrimaryIP   string
	SecondaryIP string
	Token       string
}

func main() {
	downloadISOFiles()
	generateIPXEConfig()
	generateHarvesterConfig()
}

func downloadISOFiles() {
	files := map[string]string{
		fmt.Sprintf("https://releases.rancher.com/harvester/%s/harvester-%s-amd64.iso", version, version):             fmt.Sprintf("%s/harvester-%s-amd64.iso", outputDir, version),
		fmt.Sprintf("https://releases.rancher.com/harvester/%s/harvester-%s-vmlinuz-amd64", version, version):         fmt.Sprintf("%s/harvester-%s-vmlinuz-amd64", outputDir, version),
		fmt.Sprintf("https://releases.rancher.com/harvester/%s/harvester-%s-initrd-amd64", version, version):          fmt.Sprintf("%s/harvester-%s-initrd-amd64", outputDir, version),
		fmt.Sprintf("https://releases.rancher.com/harvester/%s/harvester-%s-rootfs-amd64.squashfs", version, version): fmt.Sprintf("%s/harvester-%s-rootfs-amd64.squashfs", outputDir, version),
	}

	err := os.MkdirAll(outputDir, os.ModePerm)
	if err != nil {
		log.Fatalf("Failed to create directory: %v", err)
	}

	for url, filePath := range files {
		if _, err := os.Stat(filePath); err == nil {
			log.Printf("File already exists, skipping: %s\n", filePath)
			continue
		} else if !os.IsNotExist(err) {
			log.Fatalf("Failed to check file status: %v", err)
		}

		log.Printf("Downloading %s to %s\n", url, filePath)
		err := downloadFile(url, filePath)
		if err != nil {
			log.Fatalf("Failed to download %s: %v", url, err)
		}
		log.Printf("Successfully downloaded %s\n", url)
	}

	log.Println("All files downloaded successfully")
}

func downloadFile(url, filePath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func generateIPXEConfig() {
	// Load templates
	bootTmpl, err := loadTemplate(fmt.Sprintf("%s/%s", templateDir, "ipxe-boot.tmpl"))
	if err != nil {
		log.Fatalf("Failed to load ipxe-boot template: %v", err)
	}

	createTmpl, err := loadTemplate(fmt.Sprintf("%s/%s", templateDir, "ipxe-create.tmpl"))
	if err != nil {
		log.Fatalf("Failed to load ipxe-create template: %v", err)
	}

	joinTmpl, err := loadTemplate(fmt.Sprintf("%s/%s", templateDir, "ipxe-join.tmpl"))
	if err != nil {
		log.Fatalf("Failed to load ipxe-join template: %v", err)
	}

	// Generate iPXE files
	files := map[string]*template.Template{
		"ipxe-boot":   bootTmpl,
		"ipxe-create": createTmpl,
		"ipxe-join":   joinTmpl,
	}

	// Create the target directory if it doesn't exist
	err = os.MkdirAll(outputDir, os.ModePerm)
	if err != nil {
		log.Fatalf("Failed to create directory: %v", err)
	}

	// Generate each file
	for fileName, tmpl := range files {
		filePath := fmt.Sprintf("%s/%s", outputDir, fileName)
		log.Printf("Generating %s\n", filePath)
		err := generateFile(tmpl, filePath, config)
		if err != nil {
			log.Fatalf("Failed to generate %s: %v", filePath, err)
		}
		log.Printf("Successfully generated %s\n", filePath)
	}

	log.Println("All files generated successfully")
}

func loadTemplate(templatePath string) (*template.Template, error) {
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return nil, err
	}
	return tmpl, nil
}

func generateFile(tmpl *template.Template, filePath string, config Config) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	err = tmpl.Execute(file, config)
	if err != nil {
		return err
	}

	return nil
}

func generateHarvesterConfig() {
	// Load templates
	createTmpl, err := loadTemplate(fmt.Sprintf("%s/%s", templateDir, "config-create.yaml.tmpl"))
	if err != nil {
		log.Fatalf("Failed to load ipxe-boot template: %v", err)
	}

	joinTmpl, err := loadTemplate(fmt.Sprintf("%s/%s", templateDir, "config-join.yaml.tmpl"))
	if err != nil {
		log.Fatalf("Failed to load ipxe-create template: %v", err)
	}

	// Generate iPXE files
	files := map[string]*template.Template{
		"config-create.yaml": createTmpl,
		"config-join.yaml":   joinTmpl,
	}

	// Create the target directory if it doesn't exist
	err = os.MkdirAll(outputDir, os.ModePerm)
	if err != nil {
		log.Fatalf("Failed to create directory: %v", err)
	}

	// Generate each file
	for fileName, tmpl := range files {
		filePath := fmt.Sprintf("%s/%s", outputDir, fileName)
		log.Printf("Generating %s\n", filePath)
		err := generateFile(tmpl, filePath, config)
		if err != nil {
			log.Fatalf("Failed to generate %s: %v", filePath, err)
		}
		log.Printf("Successfully generated %s\n", filePath)
	}

	log.Println("All files generated successfully")
}
