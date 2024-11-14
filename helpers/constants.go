package helpers

import "os"

const ImaluumCasPage = "https://cas.iium.edu.my:8448/cas/login?service=https%3a%2f%2fimaluum.iium.edu.my%2fhome"

const ImaluumProfilePage = "https://imaluum.iium.edu.my/Profile"

const ImaluumLoginPage = "https://cas.iium.edu.my:8448/cas/login?service=https%3a%2f%2fimaluum.iium.edu.my%2fhome?service=https%3a%2f%2fimaluum.iium.edu.my%2fhome"

const ImaluumHomePage = "https://imaluum.iium.edu.my/home"

const ImaluumSchedulePage = "https://imaluum.iium.edu.my/MyAcademic/schedule"

const ImaluumResultPage = "https://imaluum.iium.edu.my/MyAcademic/result"

const TimeSeparator = "-"

func GetOpenAPISpecPath() string {
	OpenAPISpecPath := "./docs/swagger/swagger.json"

	if os.Getenv("APP_ENV") == "production" {
		OpenAPISpecPath = "/swagger.json"
	}
	return OpenAPISpecPath
}

