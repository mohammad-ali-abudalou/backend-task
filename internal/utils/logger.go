package utils

import "log"

// Info Logs An Informational Message :
func Info(msg string) {

	log.Printf("[INFO] %s", msg)
}

// Error Logs An Error Message :
func Error(msg string) {

	log.Printf("[ERROR] %s", msg)
}

// Fatal Logs A Fatal Message And Exits The Application :
func Fatal(msg string) {

	log.Fatalf("[FATAL] %s", msg)
}
