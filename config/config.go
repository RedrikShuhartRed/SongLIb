package config

import "os"

const (
	defaultLibraryPort     = "8081"
	defaultExternalPort    = "8080"
	defaultLibraryUser     = "postgres"
	defaultLibraryPassword = "root"
	defaultLibraryHost     = "127.0.0.1"
	defaultLibraryDBPort   = "5432"
	defaultLibrarySSLMode  = "disable"
	defaultLibraryName     = "music_library"
)

type Config struct {
	LibraryPort     string
	ExternalPort    string
	LibraryUser     string
	LibraryPassword string
	LibraryHost     string
	LibraryDBPort   string
	LibrarySSLMode  string
	LibraryName     string
}

func NewConfig() *Config {
	libraryPort := os.Getenv("LIBRARY_PORT")
	if libraryPort == "" {
		libraryPort = defaultLibraryPort
	}
	externalPort := os.Getenv("EXTERNAL_PORT")
	if externalPort == "" {
		externalPort = defaultExternalPort
	}
	libraryUser := os.Getenv("LIBRARY_USER")
	if libraryUser == "" {
		libraryUser = defaultLibraryUser
	}
	libraryPassword := os.Getenv("LIBRARY_PASSWORD")
	if libraryPassword == "" {
		libraryPassword = defaultLibraryPassword
	}
	libraryHost := os.Getenv("LIBRARY_HOST")
	if libraryHost == "" {
		libraryHost = defaultLibraryHost
	}
	libraryDBPort := os.Getenv("LIBRARY_DBPORT")
	if libraryDBPort == "" {
		libraryDBPort = defaultLibraryDBPort
	}
	librarySSLMode := os.Getenv("LIBRARY_SSLMODE")
	if librarySSLMode == "" {
		librarySSLMode = defaultLibrarySSLMode
	}

	libraryName := os.Getenv("LIBRARY_NAME")
	if libraryName == "" {
		libraryName = defaultLibraryName
	}
	return &Config{
		LibraryPort:     libraryPort,
		ExternalPort:    externalPort,
		LibraryUser:     libraryUser,
		LibraryPassword: libraryPassword,
		LibraryHost:     libraryHost,
		LibraryDBPort:   libraryDBPort,
		LibrarySSLMode:  librarySSLMode,
		LibraryName:     libraryName,
	}
}
